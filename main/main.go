package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	. "project/main/event"

	"github.com/gocolly/colly"
)

func main() {
	data := bytesFromAPI()

	events := eventCreator(data)

	for _, event := range events {
		fmt.Println(event.Type)
	}

	saveInArchive(events)

}

// From byte data to structs of type event
func eventCreator(data []byte) []Event {
	var events []Event
	if err := json.Unmarshal(data, &events); err != nil {
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
	file, err := os.OpenFile("archive/archive.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
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
		log.Fatal(err)
	}
}
