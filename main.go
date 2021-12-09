package main

import (
	"log"
	"net/http"

	"go-bootcamp/data"
	"go-bootcamp/endpoint"
)

func main() {
	repository := data.NewCsv(data.NewCsvFileBridge("pokemon.csv"))
	info := &endpoint.PokemonInfo{Pokedex: repository}
	habitat := &endpoint.Habitat{Pokedex: repository}
	all := &endpoint.All{Pokedex: repository}

	http.HandleFunc("/info", endpoint.WrapHandler(info))
	http.HandleFunc("/habitat", endpoint.WrapHandler(habitat))
	http.HandleFunc("/all", endpoint.WrapHandler(all))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
