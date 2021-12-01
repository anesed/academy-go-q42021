package data

import (
	"io"
	"strings"
	"testing"
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
	reader := fakeReader{content: []byte(strings.Join(args, "\n"))}
	bridge := fakeBridge{reader: reader}

	return bridge
}

type fakeBridge struct {
	reader fakeReader
}

func (bridge fakeBridge) openReader() (io.ReadCloser, error) {
	return &bridge.reader, nil
}

type fakeReader struct {
	content []byte
	current int
}

func (reader *fakeReader) Read(out []byte) (int, error) {
	bytes := len(out)
	remaining := len(reader.content) - reader.current

	if remaining == 0 {
		return 0, io.EOF
	}

	if remaining < bytes {
		bytes = remaining
	}

	for i := range out[:bytes] {
		out[i] = reader.content[reader.current]
		reader.current++
	}

	return bytes, nil
}

func (reader *fakeReader) Close() error {
	return nil
}
