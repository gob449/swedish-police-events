package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	. "project/main/event"
	"sort"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func main() {

	// Archive data
	eventsInArchive := getArchive()
	// New data
	newEvents := getNewEvents()
	// Merge old event with new events. Also, save the amount of duplicates in variable (could be useful)
	mergedEvents, duplicates := mergeEvents(eventsInArchive, newEvents)

	// Save new slice of events in archive
	saveInArchive(mergedEvents)
	fmt.Println("Events where successfully fetched, created and saved. \nAmount of duplicates were:", duplicates)

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
	//ADD /main BEFORE PUSH
	data, _ := os.ReadFile("archive/archive.json")
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
	//ADD /main BEFORE PUSH
	file, err := os.OpenFile("archive/archive.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("An error occurred while trying to open the archive")
		log.Fatal(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	data, err := json.Marshal(events)
	if err != nil {
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		fmt.Println("An error occurred while trying to store data in the archive")
		log.Fatal(err)
	}
}

func terminalTemplate() {
	fmt.Println("-------------------------------------------------------------------------------------------------------")
	fmt.Printf("This program gather information about crimes in Sweden, posted on the Swedish police website.\nYou can sort the crimes by different catagories.\nWrite one of the following words to sort the crimes by it:\n1. Type\n2. Location\n3. Datetime\n4. ID\nWrite 'exit' if you want to exit the program\n")
	var category string
	fmt.Scanln(&category)
	lowerCategory := strings.ToLower(category)
	switch lowerCategory {
	case "type":
		printSpecificTypeInTerminal()
	case "location":
		printSpecificLocationInTerminal()
	case "datetime":
		printDatetimesInTerminal()
	case "id":
		printIdsInTermianl()
	case "exit":
		print("Ok, exiting the program...\n")
		time.Sleep(3 * time.Second)
		os.Exit(1)
	default:
		print("Wrong input, try again\n")
		terminalTemplate()
	}
}

func printSpecificTypeInTerminal() {
	fmt.Printf("Write a specific type of event to get all crimes of that type, or write 'all' to get all crimes sorted by the types in alpabethical order\n")

	var typeSearch string
	fmt.Scanln(&typeSearch)
	lowerType := strings.ToLower(typeSearch)

	eventsInArchive := getArchive()
	sort.Slice(eventsInArchive, func(i, j int) bool {
		return eventsInArchive[i].Type < eventsInArchive[j].Type
	})
	for _, event := range eventsInArchive {
		currentType := strings.ToLower(event.Type)
		if lowerType != "all" && lowerType == currentType {
			fmt.Println(event.Name)
		} else if lowerType == "all" {
			fmt.Println(event.Name)
		}
	}
}

func printSpecificLocationInTerminal() {
	fmt.Printf("Write a specific location to get all crimes that happened in that location, or write 'all' to get all crimes sorted by the location in alpabethical order\n")

	var locationSearch string
	fmt.Scanln(&locationSearch)
	lowerLocation := strings.ToLower(locationSearch)

	eventsInArchive := getArchive()
	sort.Slice(eventsInArchive, func(i, j int) bool {
		return eventsInArchive[i].Location.Name < eventsInArchive[j].Location.Name
	})
	for _, event := range eventsInArchive {
		currentlocation := strings.ToLower(event.Location.Name)
		if lowerLocation != "all" && lowerLocation == currentlocation {
			fmt.Println(event.Name)
		} else if lowerLocation == "all" {
			fmt.Println(event.Name)
		}
	}
}

func printDatetimesInTerminal() {
	eventsInArchive := getArchive()
	sort.Slice(eventsInArchive, func(i, j int) bool {
		return eventsInArchive[i].Datetime < eventsInArchive[j].Datetime
	})
	for _, event := range eventsInArchive {
		fmt.Println(event.Name)
	}
}

func printIdsInTermianl() {
	eventsInArchive := getArchive()
	sort.Slice(eventsInArchive, func(i, j int) bool {
		return eventsInArchive[i].Id < eventsInArchive[j].Id
	})
	for _, event := range eventsInArchive {
		fmt.Println(event.Name)
	}
}
