package mock

import "io"

type fakeReader struct {
	content []byte
	current int
}

// Returns a ReadCloser object that reads the given string
func NewFakeReader(content string) io.ReadCloser {
	return &fakeReader{content: []byte(content)}
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
