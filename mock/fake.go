package mock

import (
	"bytes"
	"io"
	"net/http"
)

type fakeReader struct {
	content []byte
	current int
}

type fakeWriter struct {
	content []byte
}

type fakeResponseWriter struct {
	fakeWriter
	headers map[string][]string
}

// Returns a ReadCloser object that reads the given string
func NewFakeReader(content string) io.ReadCloser {
	return &fakeReader{content: []byte(content)}
}

// Returns a WriteCloser object that stores what it receives in memory and an Equal function to compare it with []byte
func NewFakeWriter() (io.WriteCloser, func([]byte) bool) {
	writer := fakeWriter{content: []byte{}}
	return &writer, (&writer).Equal
}

// Returns an http.ResponseWriter object
func NewFakeResponseWriter() (http.ResponseWriter, func() []byte) {
	writer := fakeResponseWriter{fakeWriter: fakeWriter{[]byte{}}, headers: map[string][]string{}}
	return &writer, (writer).GetWrittenBytes
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

func (writer *fakeWriter) Write(payload []byte) (int, error) {
	writer.content = append(writer.content, payload...)
	return len(payload), nil
}

func (writer *fakeWriter) GetWrittenBytes() []byte {
	return writer.content
}

func (writer *fakeWriter) Equal(expected []byte) bool {
	return bytes.Equal(expected, writer.content)
}

func (writer *fakeWriter) Close() error {
	return nil
}

func (writer *fakeResponseWriter) Header() http.Header {
	return writer.headers
}

func (writer *fakeResponseWriter) WriteHeader(statusCode int) {
	// Change if headers become important later
}
