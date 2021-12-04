package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-bootcamp/data"
)

type PokemonInfo struct {
	Pokedex data.Pokedex
	ID      int `query:"id,required"`
}

func (i *PokemonInfo) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	pokemon, err := i.Pokedex.Get(i.ID)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(fmt.Sprintf(`{"error":"Pokemon with ID %d not found"}`, i.ID)))

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
