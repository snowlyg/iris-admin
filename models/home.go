package models

type HomeDate struct {
	Clients   int `json:"clients"`
	Companies int `json:"companies"`
	Orders    int `json:"orders"`
}
