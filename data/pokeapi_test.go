package data

import (
	"net/http"
	"testing"

	"go-bootcamp/mock"
)

var endpoints map[string]string = map[string]string{
	"/pokemon/1": `
		{
			"id": 1,
			"species": {
				"name": "bulbasaur",
				"url": "https://mockapi/pokemon-species/100"
			}

		}`,
	"https://mockapi/pokemon-species/100": `
		{
			"id": 100,
			"name": "bulbasaur",
			"habitat": {
				"name": "grassland",
				"url": "https://mockapi/habitat/1"
			}
		}`,
}

func TestGetHabitat(t *testing.T) {
	client := http.Client{Transport: fakeTransport{}}
	bridge := NewHttpPokeapiBridge(client, "")

	habitat, err := bridge.GetHabitatFor(1)

	if err != nil {
		t.Log("bridge.GetHabitatFor should not fail")
	}

	if habitat != "grassland" {
		t.Log("Actual habitat ", habitat, " does not match expected value 'grassland'")
		t.Fail()
	}
}

type fakeTransport struct {
}

type fakeReader struct {
	Contents string
}

func (transport fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	payload, found := endpoints[req.URL.String()]

	var response http.Response

	if found {
		response = http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       mock.NewFakeReader(payload),
		}
	} else {
		response = http.Response{
			Status:     "404 Not Found",
			StatusCode: 404,
			Body:       mock.NewFakeReader("{}"),
		}
	}

	return &response, nil
}
