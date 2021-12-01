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

	http.Handle("/info", info)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
