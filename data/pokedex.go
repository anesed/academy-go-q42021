package data

import "go-bootcamp/model"

type Pokedex interface {
	Get(id int) (model.Pokemon, error)
	All() ([]model.Pokemon, error)
	Update(model.Pokemon) error
}
