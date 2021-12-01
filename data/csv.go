package data

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"go-bootcamp/model"
)

type fileBridge interface {
	openReader() (io.ReadCloser, error)
}

type Csv struct {
	index       map[int]model.Pokemon
	initialized bool
	bridge      fileBridge
}

type csvFileBridge struct {
	file string
}

func (bridge csvFileBridge) openReader() (io.ReadCloser, error) {
	return os.Open(bridge.file)
}

// Returns a new fileBridge to be consumed by a Csv
func NewCsvFileBridge(fileName string) fileBridge {
	bridge := csvFileBridge{file: fileName}

	return bridge
}

// Returns a new Csv from the provided file
func NewCsv(bridge fileBridge) Csv {
	csv := Csv{bridge: bridge}

	return csv
}

// Returns all Pokemon in storage
func (storage Csv) All() ([]model.Pokemon, error) {
	err := (&storage).init()
	if err != nil {
		return []model.Pokemon{}, err
	}

	data := make([]model.Pokemon, len(storage.index))
	i := 0

	for _, pokemon := range storage.index {
		data[i] = pokemon
	}

	return data, err
}

// Returns a single Pokemon object by its ID
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
		storage.initialized = err == nil
	}

	return err
}

func (storage *Csv) readFromFile() error {
	file, err := storage.bridge.openReader()
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

		pokemon, err := lineToPokemon(line)
		if err != nil {
			return err
		}

		storage.index[pokemon.ID] = pokemon
	}

	return nil
}

func lineToPokemon(line []string) (model.Pokemon, error) {
	if len(line) < 2 {
		return model.Pokemon{}, errors.New("Invalid record line")
	}

	id, err := strconv.Atoi(line[0])

	if err != nil {
		return model.Pokemon{}, err
	}

	pokemon := model.Pokemon{ID: id, Name: line[1]}

	return pokemon, nil
}
