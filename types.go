package dondeestas

import (
	"time"
)

type Person struct {
	Id       int    `json:"_id"` // TODO: this is very couchdb-specific?
	Name     string `json:"name"`
	Position struct {
		Tov       time.Time `json:"tov"`
		Latitude  float32   `json:"latitude"`
		Longitude float32   `json:"longitude"`
		Elevation float32   `json:"elevation"`
	} `json:"position"`
	Visible   bool  `json:"visible"`
	Whitelist []int `json:"whitelist"`
	Following []int `json:"following"`
}
