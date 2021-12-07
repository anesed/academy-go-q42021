package data

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"sync/atomic"

	"go-bootcamp/concurrency"
	"go-bootcamp/model"
)

type fileBridge interface {
	openReader() (io.ReadCloser, error)
	openWriter() (io.WriteCloser, error)
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
	return os.OpenFile(bridge.file, os.O_RDONLY, 0644)
}

func (bridge csvFileBridge) openWriter() (io.WriteCloser, error) {
	return os.OpenFile(bridge.file, os.O_WRONLY, 0644)
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
func (storage Csv) All(limit int, criteria string, workers int, itemsPerWorker int) ([]model.Pokemon, error) {
	r := 1
	var filter func(model.Pokemon) bool
	switch criteria {
	case "even":
		r = 0
		fallthrough
	case "odd":
		filter = func(p model.Pokemon) bool {
			return p.ID%2 == r
		}
	default:
		filter = func(model.Pokemon) bool {
			return true
		}
	}

	data := make([]model.Pokemon, 0)
	err := storage.readFromFile(limit, func(p model.Pokemon) {
		data = append(data, p)
	}, filter, concurrency.NewWorkerPool(workers, itemsPerWorker, limit))

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

// Saves an updated pokemon to the data store
func (storage Csv) Update(pokemon model.Pokemon) error {
	err := (&storage).init()
	if err != nil {
		return err
	}

	storage.index[pokemon.ID] = pokemon
	contents := make([][]string, len(storage.index))
	index := 0

	for _, line := range storage.index {
		contents[index] = pokemonToLine(line)
		index++
	}

	file, err := storage.bridge.openWriter()
	if err != nil {
		return err
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.WriteAll(contents)

	return err
}

func (storage *Csv) init() error {
	var err error = nil

	if !storage.initialized {
		storage.index = make(map[int]model.Pokemon)
		err := storage.readFromFile(0, func(p model.Pokemon) {
			storage.index[p.ID] = p
		}, func(model.Pokemon) bool { return true }, nil)
		storage.initialized = err == nil
	}

	return err
}

func (storage *Csv) readFromFile(limit int, store func(model.Pokemon), filter func(model.Pokemon) bool, wp *concurrency.WorkerPool) error {
	file, err := storage.bridge.openReader()
	if err != nil {
		return err
	}

	if wp == nil {
		wp = concurrency.NewWorkerPool(10, 0, 0)
	}
	defer wp.Close()
	defer file.Close()

	reader := csv.NewReader(file)
	output := make(chan model.Pokemon)
	processed := int64(0)

	go func() {
		for entry := range output {
			if limit == 0 || processed < int64(limit) {
				store(entry)
				atomic.AddInt64(&processed, 1)
			}
		}
	}()

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if limit > 0 && atomic.LoadInt64(&processed) > int64(limit) {
			close(output)
			break
		}

		wp.Push(func() bool {
			pokemon, _ := lineToPokemon(line)
			if filter(pokemon) {
				output <- pokemon
				return true
			}
			return false
		})
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
	if len(line) > 2 {
		pokemon.Habitat = line[2]
	}

	return pokemon, nil
}

func pokemonToLine(pokemon model.Pokemon) []string {
	return []string{strconv.Itoa(pokemon.ID), pokemon.Name, pokemon.Habitat}
}
