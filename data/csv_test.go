package data

import (
	"io"
	"strings"
	"testing"

	"go-bootcamp/mock"
)

func TestCsvGetAll(t *testing.T) {
	bridge := buildFakeBridge("1,bulbasaur", "2,ivysaur", "3,venusaur")
	csv := NewCsv(bridge)

	records, _ := csv.All()
	expectedCount := 3

	if len(records) != expectedCount {
		t.Logf("Expected %d records, got %d", expectedCount, len(records))
		t.Fail()
	}
}

func TestCsvGetOne(t *testing.T) {
	bridge := buildFakeBridge("1,bulbasaur")
	csv := NewCsv(bridge)

	record, _ := csv.Get(1)

	if record.ID != 1 || record.Name != "bulbasaur" {
		t.Log("Invalid record", record)
		t.Fail()
	}

	record, err := csv.Get(2)

	if err != nil {
		t.Log("Error expected but none received")
	}
}

func buildFakeBridge(args ...string) fakeBridge {
	reader := mock.NewFakeReader(strings.Join(args, "\n"))
	bridge := fakeBridge{reader: reader}

	return bridge
}

type fakeBridge struct {
	reader io.ReadCloser
}

func (bridge fakeBridge) openReader() (io.ReadCloser, error) {
	return bridge.reader, nil
}
