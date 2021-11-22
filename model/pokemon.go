package model

type Pokemon struct {
	ID   int    `json:"pokedex_entry"`
	Name string `json:"name"`
}

type Pokedex struct {
	entries map[int]Pokemon
}

func (dex *Pokedex) Get(id int) (Pokemon, bool) {
	entry, found := dex.entries[id]

	return entry, found
}

func NewDex(pokemon []Pokemon) Pokedex {
	data := map[int]Pokemon{}
	dex := Pokedex{entries: data}
	for _, entry := range pokemon {
		dex.entries[entry.ID] = entry
	}

	return dex
}
