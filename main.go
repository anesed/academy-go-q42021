package main

import (
	"go-bootcamp/action"
	"go-bootcamp/data"
	"go-bootcamp/model"
	"log"
	"net/http"
)

func main() {
	storage := data.Csv{}
	dex := model.NewDex(storage.All())
	info := &action.PokemonInfo{Pokedex: dex}
	http.Handle("/info", info)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
