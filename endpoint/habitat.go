package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-bootcamp/data"
)

type Habitat struct {
	Pokedex data.Pokedex
	ID      int `query:"id,required"`
}

func (i *Habitat) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	pokemon, err := i.Pokedex.Get(i.ID)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf(`{"error":"Pokemon with ID %d not found"}`, i.ID)))

		return
	}

	if len(pokemon.Habitat) == 0 {
		habitat, err := fetchPokemonHabitat(pokemon.ID)
		if err == nil {
			pokemon.Habitat = habitat
			i.Pokedex.Update(pokemon)
		}
	}

	writer.WriteHeader(http.StatusOK)
	body, err := json.Marshal(pokemon)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"error":"Internal server error"}`))

		return
	}

	writer.Write(body)
}

func fetchPokemonHabitat(pokemonId int) (string, error) {
	api := data.NewHttpPokeapiBridge(http.Client{}, "https://pokeapi.co/api/v2")

	return api.GetHabitatFor(pokemonId)
}
