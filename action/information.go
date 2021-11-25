package action

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"go-bootcamp/data"
)

type PokemonInfo struct {
	Pokedex data.Pokedex
}

func (i *PokemonInfo) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
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

	writer.WriteHeader(http.StatusOK)
	body, err := json.Marshal(pokemon)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"error":"Internal server error"}`))

		return
	}

	writer.Write(body)
}
