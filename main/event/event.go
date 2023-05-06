// Gabriel Räätäri Nyström 2023-04-26
// Hugo Larsson Wilhelmsson 2023-04-26
// This class defines the data structure 'event' and provides
// sorting implmentations for 'event' based on different criteria
package event

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/pkg/browser"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	TypeKeys = []string{
		"Alkohollagen",
		"Anträffad död",
		"Anträffat gods",
		"Arbetsplatsolycka",
		"Bedrägeri",
		"Bombhot",
		"Brand",
		"Brand automatlarm",
		"Bråk",
		"Detonation",
		"Djur skadat/omhändertaget",
		"Ekobrott",
		"Farligt föremål, misstänkt",
		"Fjällräddning",
		"Fylleri/LOB",
		"Förfalskningsbrott",
		"Försvunnen person",
		"Gränskontroll",
		"Häleri",
		"Inbrott",
		"Inbrott, försök",
		"Knivlagen",
		"Kontroll person/fordon",
		"Lagen om hundar och katter",
		"Larm inbrott",
		"Larm överfall",
		"Miljöbrott",
		"Missbruk av urkund",
		"Misshandel",
		"Misshandel, grov",
		"Mord/dråp",
		"Mord/dråp, försök",
		"Motorfordon, anträffat stulet",
		"Motorfordon, stöld",
		"Narkotikabrott",
		"Naturkatastrof",
		"Ofog barn/ungdom",
		"Ofredande/förargelse",
		"Olaga frihetsberövande",
		"Olaga hot",
		"Olaga intrång/hemfridsbrott",
		"Olovlig körning",
		"Ordningslagen",
		"Polisinsats/kommendering",
		"Rattfylleri",
		"Rån",
		"Rån väpnat",
		"Rån övrigt",
		"Rån, försök",
		"Räddningsinsats",
		"Sammanfattning dag",
		"Sammanfattning dygn",
		"Sammanfattning eftermiddag",
		"Sammanfattning förmiddag",
		"Sammanfattning helg",
		"Sammanfattning kväll",
		"Sammanfattning kväll och natt",
		"Sammanfattning natt",
		"Sammanfattning vecka",
		"Sedlighetsbrott",
		"Sjukdom/olycksfall",
		"Sjölagen",
		"Skadegörelse",
		"Skottlossning",
		"Skottlossning, misstänkt",
		"Spridning smittsamma kemikalier",
		"Stöld",
		"Stöld, försök",
		"Stöld, ringa",
		"Stöld/inbrott",
		"Tillfälligt obemannat",
		"Trafikbrott",
		"Trafikhinder",
		"Trafikkontroll",
		"Trafikolycka",
		"Trafikolycka, personskada",
		"Trafikolycka, singel",
		"Trafikolycka, smitning från",
		"Trafikolycka, vilt",
		"Uppdatering",
		"Utlänningslagen",
		"Vapenlagen",
		"Varningslarm/haveri",
		"Våld/hot mot tjänsteman",
		"Våldtäkt",
		"Våldtäkt, försök",
		"Vållande till kroppsskada",
	}
)

type Event struct {
	Id       int    `json:"id"`
	Datetime string `json:"datetime"`
	Name     string `json:"name"`
	Summary  string `json:"summary"`
	Url      string `json:"url"`
	Type     string `json:"type"`
	Location struct {
		Name string `json:"name"`
		Gps  string `json:"gps"`
	} `json:"location"`
}

// MarshalJSON Implements []byte() method
func (e *Event) MarshalJSON() ([]byte, error) {
	type Alias Event
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	})
}

// ById Implements sort methods:
// Len() int Less(i int, j int) bool Swap(i int, j int)
type ById []Event

func (e ById) Len() int {
	return len(e)
}

func (e ById) Less(i, j int) bool {
	return e[i].Id < e[j].Id
}

func (e ById) Swap(i int, j int) {
	e[i], e[j] = e[j], e[i]
}

// ByDatetime Implements sort methods:
// Len() int Less(i int, j int) bool Swap(i int, j int)
type ByDatetime []Event

func (e ByDatetime) Len() int {
	return len(e)
}

func (e ByDatetime) Less(i, j int) bool {
	t1, err := time.Parse("2006-01-02 15:04:05 -07:00", e[i].Datetime)
	if err != nil {
		fmt.Println("Failed to parse", e[i].Datetime)
		log.Fatal(err)
	}
	t2, err := time.Parse("2006-01-02 15:04:05 -07:00", e[j].Datetime)
	if err != nil {
		fmt.Println("Failed to parse", e[j].Datetime)
		log.Fatal(err)
	}
	return t1.Before(t2)
}

func (e ByDatetime) Swap(i int, j int) {
	e[i], e[j] = e[j], e[i]
}

// ByType Implements sort methods:
// Len() int Less(i int, j int) bool Swap(i int, j int)
type ByType []Event

func (e ByType) Len() int {
	return len(e)
}

func (e ByType) Less(i, j int) bool {
	return strings.Compare(e[i].Type, e[j].Type) < 0
}

func (e ByType) Swap(i int, j int) {
	e[i], e[j] = e[j], e[i]
}

// ByLocation Implements sort methods:
// Len() int Less(i int, j int) bool Swap(i int, j int)
type ByLocation []Event

func (e ByLocation) Len() int {
	return len(e)
}

func (e ByLocation) Less(i, j int) bool {
	return strings.Compare(e[i].Location.Name, e[j].Location.Name) < 0
}

func (e ByLocation) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

// AllEventsSlice Merges archive- and new data and returns the merged slice
func AllEventsSlice() []Event {
	// Archive data
	eventsInArchive := GetArchive()
	// New data
	newEvents := GetNewEvents()
	// Merge old event with new events. Also, save the amount of duplicates in variable (could be useful)
	mergedEvents, _ := MergeEvents(eventsInArchive, newEvents)

	sort.Sort(ByDatetime(mergedEvents))

	return mergedEvents
}

func SubCatType(events []Event, key string) []Event {
	var subCategory []Event

	for _, event := range events {
		if event.Type == key {
			subCategory = append(subCategory, event)
		}
	}

	return subCategory
}

func SubCatLocation(events []Event, key string) []Event {
	var subCategory []Event

	for _, event := range events {
		if event.Location.Name == key {
			subCategory = append(subCategory, event)
		}
	}

	return subCategory
}

// Merges new and old events into a single slice
func MergeEvents(eventsInArchive []Event, newEvents []Event) ([]Event, int) {
	// Old and new with duplicates
	allEventsRaw := append(eventsInArchive, newEvents...)
	var mergedEvents []Event
	visited := make(map[int]int)
	var duplicates int
	// Time complexity is O(n^2) because "append" creates a new slice which is O(n) for every new event
	for _, event := range allEventsRaw {
		_, ok := visited[event.Id]
		if !ok {
			visited[event.Id]++
			mergedEvents = append(mergedEvents, event)
		} else {
			duplicates++
		}
	}
	return mergedEvents, duplicates
}

// returns the new data of type event
func GetNewEvents() []Event {
	newData := bytesFromAPI()
	newEvents := eventCreator(newData)
	return newEvents
}

// Returns the events that are currently stored in the archive
func GetArchive() []Event {
	data, _ := os.ReadFile("main/archive/archive.json")
	eventsInArchive := eventCreator(data)
	return eventsInArchive
}

// From byte data to structs of type event
func eventCreator(data []byte) []Event {
	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
		fmt.Println("Error occurred while creating events")
		log.Fatal(err)
	}
	return events
}

// Scrapes API-data and returns byte arr
func bytesFromAPI() []byte {
	var data []byte
	c := colly.NewCollector()
	c.OnError(func(response *colly.Response, err error) {
		log.Fatal(err)
	})
	c.OnResponse(func(response *colly.Response) {
		data = response.Body
	})
	if err := c.Visit("https://polisen.se/api/events"); err != nil {
		log.Fatal(err)
	}
	return data
}

// SaveInArchive saves a slice of events in a JSON file located at "main/archive/archive.json".
// If the file doesn't exist, it creates the file. If the file already exists, it appends the data.
// The function returns an error if any error occurs during file operations.
func SaveInArchive(events []Event) {
	// Creates file if necessary, appends if file exists
	file, err := os.OpenFile("main/archive/archive.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("An error occurred while trying to open the archive")
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("An error occurred while trying to close the archive")
			log.Fatal(err)
		}
	}()
	writer := bufio.NewWriter(file)
	data, err := json.Marshal(events)
	if err != nil {
		fmt.Println("An error occurred while trying to store data in the archive")
	}
	if _, err := writer.Write(data); err != nil {
		fmt.Println("An error occurred while trying to store data in the archive")
		log.Fatal(err)
	}
}

// Takes Event.URL value and opens a webpage with the corresponding extensive event summary
func OpenInBrowser(URL string) {
	url := "https://polisen.se/" + URL
	if err := browser.OpenURL(url); err != nil {
		fmt.Println("An error occurred while trying to open the event summary page.")
	}
}

func GetExtendedSummary(URL string) string {
	url := "https://polisen.se/" + URL

	var summary string

	c := colly.NewCollector(colly.AllowURLRevisit())

	c.OnError(func(response *colly.Response, err error) {
		log.Fatal(err)
	})

	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101")
	})

	c.OnHTML("#main-content > div.body-content-wrapper > div > div > div > div > div.event-content > div.text-body.editorial-html", func(e *colly.HTMLElement) {
		summary = strings.TrimSpace(e.Text)
	})

	if err := c.Visit(url); err != nil {
		log.Fatal(err)
	}

	return summary
}
