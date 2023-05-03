// Gabriel Räätäri Nyström 2023-04-26
// Hugo Larsson Wilhelmsson 2023-04-26
// This class defines the data structure 'event' and provides
// sorting implmentations for 'event' based on different criteria
package event

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
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
