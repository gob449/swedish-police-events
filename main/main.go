package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

func main() {
	data := bytesFromAPI()
	fmt.Println(data)
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
