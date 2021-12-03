package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-bootcamp/data"
)

type Habitat struct {
	Pokedex data.Pokedex
}

func (i *Habitat) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	id, found := r.URL.Query()["id"]
	if !found {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"error":"ID not provided"}`))

		return
	}

	parsedId, err := strconv.Atoi(id[0])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf(`{"error":"Invalid ID: %v"}`, id[0])))

		return
	}

	pokemon, err := i.Pokedex.Get(parsedId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf(`{"error":"Pokemon with ID %d not found"}`, parsedId)))

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
