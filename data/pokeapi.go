package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type pokeapiBridge interface {
	GetHabitatFor(pokemonId int) (string, error)
}

type httpPokeapiBridge struct {
	client  http.Client
	baseUrl string
}

// Creates a new bridge to connect with pokeapi
func NewHttpPokeapiBridge(client http.Client, baseUrl string) pokeapiBridge {
	return httpPokeapiBridge{client: client, baseUrl: baseUrl}
}

// Returns the habitat of the pokemon with the given ID
func (bridge httpPokeapiBridge) GetHabitatFor(pokemonId int) (string, error) {
	speciesUrl, err := bridge.fetchSpeciesUrl(pokemonId)

	if err != nil {
		log.Println(err)
		return "", errors.New("Api fetch error")
	}

	var speciesData struct {
		Habitat struct {
			Name string `json:"name"`
		} `json:"habitat"`
	}

	err = bridge.fetchResourceData(speciesUrl, &speciesData)

	if err != nil {
		return "", errors.New("Error decoding API response")
	}

	return speciesData.Habitat.Name, nil
}

func (bridge httpPokeapiBridge) fetchSpeciesUrl(pokemonId int) (string, error) {
	var pokemonData struct {
		Species struct {
			Url string `json:"url"`
		} `json:"species"`
	}

	err := bridge.fetchResourceData(bridge.resourceUrl("pokemon", pokemonId), &pokemonData)

	if err != nil {
		return "", err
	}

	return pokemonData.Species.Url, nil
}

func (bridge httpPokeapiBridge) resourceUrl(resource string, resourceId int) string {
	return fmt.Sprintf(fmt.Sprintf("%s/%s/%d", bridge.baseUrl, resource, resourceId))
}

func (bridge httpPokeapiBridge) fetchResourceData(url string, output interface{}) error {
	res, err := bridge.client.Get(url)

	if err != nil {
		return fmt.Errorf("Connection error: %w", err)
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&output)

	if err != nil {
		return fmt.Errorf("Decoding error: %w", err)
	}

	return nil
}
