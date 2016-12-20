package dondeestas

import "time"

// Person encapsulates all of the data needed for tracking a user
type Person struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position struct {
		Tov       time.Time `json:"tov"`
		Latitude  float32   `json:"latitude"`
		Longitude float32   `json:"longitude"`
		Elevation float32   `json:"elevation"`
	} `json:"position"`
	Visible   bool     `json:"visible"`
	Whitelist []string `json:"whitelist"`
	Following []string `json:"following"`
}
