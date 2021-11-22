package model

type Pokemon struct {
	ID   int    `json:"pokedex_entry"`
	Name string `json:"name"`
}
