package main

import (
	"log"
	"net/http"

	"go-bootcamp/action"
	"go-bootcamp/data"
)

func main() {
	repository := data.Csv{}
	info := &action.PokemonInfo{Pokedex: repository}
	http.Handle("/info", info)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
