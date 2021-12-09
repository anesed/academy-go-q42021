package data

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"go-bootcamp/mock"
)

func TestCsvGetAll(t *testing.T) {
	bridge := buildFakeBridge("1,bulbasaur", "2,ivysaur", "3,venusaur")
	csv := NewCsv(bridge)

	records, _ := csv.All(0, "", 1, 0)
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

	if err == nil {
		t.Log("Error expected but none received")
		t.Fail()
	}
}

func TestCsvErrorInAll(t *testing.T) {
	bridge := buildFakeBridge("1,bulbasaur", "2", "3,venusaur")
	csv := NewCsv(bridge)

	_, err := csv.All(0, "", 1, 0)

	if err == nil {
		t.Log("Error expected but none received")
		t.Fail()
	}
}

func TestCsvUpdate(t *testing.T) {
	bridge := buildFakeBridge("1,bulbasaur")
	writer, getContents := mock.NewFakeWriter()
	bridge.writer = writer
	csv := NewCsv(bridge)

	record, _ := csv.Get(1)
	record.Habitat = "grassland"

	csv.Update(record)
	expected := "1,bulbasaur,grassland\n"

	if !bytes.Equal([]byte(expected), getContents()) {
		t.Logf("Written content `%s`doesn't match expected %s", string(getContents()), expected)
		t.Fail()
	}
}

func buildFakeBridge(args ...string) fakeBridge {
	reader := mock.NewFakeReader(strings.Join(args, "\n"))
	bridge := fakeBridge{reader: reader}

	return bridge
}

type fakeBridge struct {
	reader io.ReadCloser
	writer io.WriteCloser
}

func (bridge fakeBridge) openReader() (io.ReadCloser, error) {
	return bridge.reader, nil
}

func (bridge fakeBridge) openWriter() (io.WriteCloser, error) {
	return bridge.writer, nil
}
