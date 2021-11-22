package data

import (
	"encoding/csv"
	"errors"
	"go-bootcamp/model"
	"io"
	"os"
	"strconv"
)

type Csv struct {
	index       map[int]model.Pokemon
	initialized bool
}

func (storage Csv) All() ([]model.Pokemon, error) {
	err := (&storage).init()
	data := make([]model.Pokemon, len(storage.index))

	for _, pokemon := range storage.index {
		data = append(data, pokemon)
	}

	return data, err
}

func (storage Csv) Get(id int) (model.Pokemon, error) {
	err := (&storage).init()
	if err != nil {
		return model.Pokemon{}, err
	}

	record, found := storage.index[id]
	if !found {
		err = errors.New("Record not found")
	}

	return record, err
}

func (storage *Csv) init() error {
	var err error = nil

	if !storage.initialized {
		err := storage.readFromFile()
		storage.initialized = err != nil
	}

	return err
}

func (storage *Csv) readFromFile() error {
	file, err := os.Open("pokemon.csv")
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	storage.index = make(map[int]model.Pokemon)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(line) != 2 {
			continue
		}

		id, err := strconv.Atoi(line[0])
		if err != nil {
			return err
		}
		pokemon := model.Pokemon{ID: id, Name: line[1]}
		storage.index[id] = pokemon
	}

	return nil
}
