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
)

func main() {
	// Archive data
	eventsInArchive := getArchive()
	// New data
	newEvents := getNewEvents()
	// Merge old event with new events. Also, save the amount of duplicates in variable (could be useful)
	mergedEvents, _ := mergeEvents(eventsInArchive, newEvents)

	openSummary(mergedEvents[0].Url)

	saveInArchive(mergedEvents)
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
		return
	}
	_, err = writer.Write(data)
	if err != nil {
		fmt.Println("An error occurred while trying to store data in the archive")
		log.Fatal(err)
	}
}

// Takes Event.Id value and opens a webpage with the corresponding extensive event summary
func openSummary(URL string) {
	URL = "https://polisen.se/" + URL
	if err := browser.OpenURL(URL); err != nil {
		fmt.Println("An error occurred while trying to open the event summary page.")
	}
}
