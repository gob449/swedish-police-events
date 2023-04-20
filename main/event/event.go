package event

import "encoding/json"

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
