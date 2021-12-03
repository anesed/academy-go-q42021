package main

import (
	"log"
	"net/http"

	"go-bootcamp/action"
	"go-bootcamp/data"
)

func main() {
	repository := data.NewCsv(data.NewCsvFileBridge("pokemon.csv"))
	info := &action.PokemonInfo{Pokedex: repository}
	habitat := &action.Habitat{Pokedex: repository}

	http.Handle("/info", info)
	http.Handle("/habitat", habitat)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
