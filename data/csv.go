package data

import (
	"encoding/csv"
	"go-bootcamp/model"
	"io"
	"log"
	"os"
	"strconv"
)

type Csv struct {
	pokedex *model.Pokedex
}

func (storage *Csv) All() []model.Pokemon {
	file, err := os.Open("pokemon.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)

	data := make([]model.Pokemon, 1)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if len(line) != 2 {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}
		pokemon := model.Pokemon{ID: id, Name: line[1]}
		data = append(data, pokemon)
	}

	return data
}
