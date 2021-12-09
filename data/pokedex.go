package data

import "go-bootcamp/model"

type Pokedex interface {
	Get(id int) (model.Pokemon, error)
	All(limit int, criteria string, workers int, itemsPerWorker int) ([]model.Pokemon, error)
	Update(model.Pokemon) error
}
