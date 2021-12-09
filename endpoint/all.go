package endpoint

import (
	"encoding/json"
	"log"
	"net/http"

	"go-bootcamp/data"
)

type All struct {
	Pokedex        data.Pokedex
	Type           string `query:"type"`
	ItemCount      int    `query:"items"`
	ItemsPerWorker int    `query:"items_per_workers"`
}

func (i *All) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	if !(i.Type == "" || i.Type == "odd" || i.Type == "even") {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"errors":["Invalid value for field type"]}`))

		return
	}

	items, err := i.Pokedex.All(i.ItemCount, i.Type, 10, i.ItemsPerWorker)
	var body []byte

	if err == nil {
		body, err = json.Marshal(items)
	}

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(`{"error":"Internal server error"}`))
		log.Println(err)

		return
	}

	if len(body) == 0 {
		body = []byte("[]")
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(body)
}
