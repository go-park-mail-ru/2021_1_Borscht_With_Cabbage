package models

type Address struct {
	Name      string `json:"name"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	Radius    int    `json:"radius"`
}

type Coordinates struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}
