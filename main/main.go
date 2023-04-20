package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	. "project/main/event"
)

func main() {
	data := bytesFromAPI()

	events := eventCreator(data)

	for _, event := range events {
		fmt.Println(event.Type)
	}

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
