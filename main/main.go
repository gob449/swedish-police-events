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
	"fmt"
	"log"
	"os"
	. "project/main/event"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Collects the new data and merges it with the previous data in the database
// Calls terminalTemplate() to start the program
func main() {

	// Run terminal program
	//terminalTemplate()

	//RunGUI()

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
	eventsInArchive := GetArchive()

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

	eventsInArchive := GetArchive()

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
	eventsInArchive := GetArchive()
	sort.Sort(ByDatetime(eventsInArchive))
	for _, event := range eventsInArchive {
		fmt.Println(event.Id, "----", event.Name)
	}
}

// Prints the names and Id:s of the crimes, based on its Id
func printIdsInTerminal() {
	eventsInArchive := GetArchive()
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
	eventsInArchive := GetArchive()
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
	eventsInArchive := GetArchive()
	for _, event := range eventsInArchive {
		if key == event.Id {
			OpenSummary(event.Url)
		}
	}
}
