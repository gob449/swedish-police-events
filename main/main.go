// Gabriel Räätäri Nyström 2023-04-26
// Hugo Larsson Wilhelmsson 2023-04-26
// This program gather information about crimes in Sweden, posten on the swedish police website.
// The crimes can be sorted by location, type, id or datetime and the user can search for
// keywords to find for example crimes connected to a specific city. The user can also get more
// detailed information about the crimes, both in the program and by opening an URL to get the whole
// description directly from the police website.
// The terminal is used to print the information and is where the user write the requests.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	. "project/main/event"
	"sort"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/gocolly/colly"
	"github.com/pkg/browser"
)

// Collects the new data and merges it with the previous data in the database
// Calls terminalTemplate() to start the program
func main() {
	/*
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
	*/

	myApp := app.New()
	myWindow := myApp.NewWindow("Hello")
	hello := widget.NewLabel("Hello Fyne!")
	myWindow.SetContent(hello)
	w2 := myApp.NewWindow("Larger")
	w2.SetContent(widget.NewLabel("Consider it configured"))
	w2.ShowAndRun()

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

// returns the new data of type event
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

// saveInArchive saves a slice of events in a JSON file located at "main/archive/archive.json".
// If the file doesn't exist, it creates the file. If the file already exists, it appends the data.
// The function returns an error if any error occurs during file operations.
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

// Prints the overall menu in the terminal
func terminalTemplate() {
	initialInfo()
	for {
		provideAlternatives()
		parseAlternativeAndAct()
	}
}

// Prints the first message in the program that tells the user what the program does
func initialInfo() {
	fmt.Printf(`
--------------------------------------------------------------------------------------------
This program gather information about crimes in Sweden, posted on the Swedish police website.
You can sort the crimes by different catagories.	
`)
}

// Prints the choices the user has in the menu
func provideAlternatives() {
	fmt.Printf(`
Write one of the following characters to sort the crimes by it:
1. t (Type) 
2. l (Location)
3. d (Datetime)
4. i (ID)
5. s (Summary)
6. v (Visit page for extensive summary)
Write 'exit' if you want to exit the program
`)
}

// Calls a specific print function depending on the users input, exit the
// program if the user wants, or tells the user if it gives the wrong input.
func parseAlternativeAndAct() {
	var category string
	if _, err := fmt.Scanln(&category); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	lowerCategory := strings.ToLower(category)
	switch lowerCategory {
	case "t":
		printSpecificTypeInTerminal()
	case "l":
		printSpecificLocationInTerminal()
	case "d":
		printDatetimeInTerminal()
	case "i":
		printIdsInTerminal()
	case "s":
		printSpecificSummaryInTerminal()
	case "v":
		parseInputForID()
	case "exit":
		print("Ok, exiting the program...\n")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	default:
		print("Wrong input, try again\n")
		terminalTemplate()
	}
}

// Prints the names and Id:s of the crimes, based on its Type.
// Allows the user to search for a specific type to get all the crimes of that type printed
func printSpecificTypeInTerminal() {
	fmt.Printf("Write a specific type of event to get all crimes of that type, or write 'all' to get all crimes sorted by the types in alpabethical order\n")

	var typeSearch string
	if _, err := fmt.Scanln(&typeSearch); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	lowerType := strings.ToLower(typeSearch)
	eventsInArchive := getArchive()

	sort.Sort(ByType(eventsInArchive))

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

// Prints the names and Id:s of the crimes, based on its Location.
// Allows the user to search for a specific location to get all the crimes connected to that location printed
func printSpecificLocationInTerminal() {
	fmt.Printf("Write a specific location to get all crimes that happened in that location, or write 'all' to get all crimes sorted by the location in alpabethical order\n")

	var locationSearch string
	if _, err := fmt.Scanln(&locationSearch); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	lowerLocation := strings.ToLower(locationSearch)

	eventsInArchive := getArchive()

	sort.Sort(ByLocation(eventsInArchive))

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

// Prints the names and Id:s of the crimes, based on its Datetime
func printDatetimeInTerminal() {
	eventsInArchive := getArchive()
	sort.Sort(ByDatetime(eventsInArchive))
	for _, event := range eventsInArchive {
		fmt.Println(event.Id, "----", event.Name)
	}
}

// Prints the names and Id:s of the crimes, based on its Id
func printIdsInTerminal() {
	eventsInArchive := getArchive()
	sort.Sort(ById(eventsInArchive))
	for _, event := range eventsInArchive {
		fmt.Println(event.Id, "----", event.Name)
	}
}

// Prints a summary of a crime connected to an Id that the user provides
func printSpecificSummaryInTerminal() {
	fmt.Println("Please provide id of the event you want to access ")
	var id string
	if _, err := fmt.Scanln(&id); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	key, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	eventsInArchive := getArchive()
	for _, event := range eventsInArchive {
		if key == event.Id {
			fmt.Println(event.Name, "----", event.Summary)
		}
	}
}

// Open an external webside with an extensive summary of a crime connected to an Id that the user provides
func parseInputForID() {
	fmt.Println("Please provide id of the event you want to access ")
	var id string
	if _, err := fmt.Scanln(&id); err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	key, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("An error occurred while parsing user input")
		log.Fatal(err)
	}
	eventsInArchive := getArchive()
	for _, event := range eventsInArchive {
		if key == event.Id {
			openSummary(event.Url)
		}
	}
}

// Takes Event.URL value and opens a webpage with the corresponding extensive event summary
func openSummary(URL string) {
	URL = "https://polisen.se/" + URL
	if err := browser.OpenURL(URL); err != nil {
		fmt.Println("An error occurred while trying to open the event summary page.")
	}
}
