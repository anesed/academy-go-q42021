package main

import (
	"go-bootcamp/action"
	"go-bootcamp/data"
	"log"
	"net/http"
)

func main() {
	repository := data.Csv{}
	info := &action.PokemonInfo{Pokedex: repository}
	http.Handle("/info", info)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
