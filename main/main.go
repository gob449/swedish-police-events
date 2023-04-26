package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/pkg/browser"
	"log"
	"os"
	. "project/main/event"
	"strings"
	"time"
)

func main() {
	// Archive data
	eventsInArchive := getArchive()
	// New data
	newEvents := getNewEvents()
	// Merge old event with new events. Also, save the amount of duplicates in variable (could be useful)
	mergedEvents, _ := mergeEvents(eventsInArchive, newEvents)
	// Save new slice of events in archive
	saveInArchive(mergedEvents)

	// start program
	terminalTemplate()

}

// Merges new and old events into a single slice
func mergeEvents(eventsInArchive []Event, newEvents []Event) ([]Event, int) {
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

func getNewEvents() []Event {
	newData := bytesFromAPI()
	newEvents := eventCreator(newData)
	return newEvents
}

// Returns the events that are currently stored in the archive
func getArchive() []Event {
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

func saveInArchive(events []Event) {
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

func terminalTemplate() {
	initialInfo()
	for {
		provideAlternatives()
		parseAlternativeAndAct()
	}
}

func parseAlternativeAndAct() {
	var category string
	if _, err := fmt.Scanln(&category); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	lowerCategory := strings.ToLower(category)
	switch lowerCategory {
	case "type":
		printSpecificTypeInTerminal()
	case "location":
		printSpecificLocationInTerminal()
	case "datetime":
		printDatetimeInTerminal()
	case "id":
		printIdsInTerminal()
	case "exit":
		print("Ok, exiting the program...\n")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	default:
		print("Wrong input, try again\n")
		terminalTemplate()
	}
}

func initialInfo() {
	fmt.Printf(`
--------------------------------------------------------------------------------------------
This program gather information about crimes in Sweden, posted on the Swedish police website.
You can sort the crimes by different catagories.	
`)
}

func provideAlternatives() {
	fmt.Printf(`
Write one of the following words to sort the crimes by it:
1. Type
2. Location
3. Datetime
4. ID
Write 'exit' if you want to exit the program
`)
}

func printSpecificTypeInTerminal() {
	fmt.Printf("Write a specific type of event to get all crimes of that type, or write 'all' to get all crimes sorted by the types in alpabethical order\n")

	var typeSearch string
	if _, err := fmt.Scanln(&typeSearch); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	lowerType := strings.ToLower(typeSearch)
	eventsInArchive := getArchive()

	if lowerType == "all" {
		for _, event := range eventsInArchive {
			fmt.Println(event.Id, "----", event.Name)
		}
	} else {
		for _, event := range eventsInArchive {
			if lowerType == strings.ToLower(event.Type) {
				fmt.Println(event.Id, "----", event.Name)
			}
		}
	}
}

func printSpecificLocationInTerminal() {
	fmt.Printf("Write a specific location to get all crimes that happened in that location, or write 'all' to get all crimes sorted by the location in alpabethical order\n")

	var locationSearch string
	if _, err := fmt.Scanln(&locationSearch); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	lowerLocation := strings.ToLower(locationSearch)

	eventsInArchive := getArchive()

	if lowerLocation == "all" {
		for _, event := range eventsInArchive {
			fmt.Println(event.Id, "----", event.Name)
		}
	} else {
		for _, event := range eventsInArchive {
			if lowerLocation == strings.ToLower(event.Location.Name) {
				fmt.Println(event.Id, "----", event.Name)
			}
		}
	}
}

func printDatetimeInTerminal() {
	eventsInArchive := getArchive()
	for _, event := range eventsInArchive {
		fmt.Println(event.Id, "----", event.Name)
	}
}

func printIdsInTerminal() {
	eventsInArchive := getArchive()
	for _, event := range eventsInArchive {
		fmt.Println(event.Id, "----", event.Name)
	}
}

// Takes Event.URL value and opens a webpage with the corresponding extensive event summary
func openSummary(URL string) {
	URL = "https://polisen.se/" + URL
	if err := browser.OpenURL(URL); err != nil {
		fmt.Println("An error occurred while trying to open the event summary page.")
	}
}
