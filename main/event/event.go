package event

import (
	"encoding/json"
	"strings"
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
// /*
type ById []Event

func (e ById) Len() int {
	return len(e)
}

func (e ById) Less(i, j int) bool {

	return e[i].Id < e[j].Id
}

func (e ById) Swap(i int, j int) {

	e[i].Id, e[j].Id = e[j].Id, e[i].Id
}

// ByDatetime Implements sort methods:
// Len() int Less(i int, j int) bool Swap(i int, j int)
type ByDatetime []Event

func (e ByDatetime) Len() int {
	return len(e)
}

func (e ByDatetime) Less(i, j int) bool {
	return e[i].Datetime < e[j].Datetime
}

func (e ByDatetime) Swap(i int, j int) {
	e[i].Datetime, e[j].Datetime = e[j].Datetime, e[i].Datetime
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
	e[i].Type, e[j].Type = e[j].Type, e[i].Type
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
