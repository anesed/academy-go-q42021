package endpoint

import (
	"bytes"
	"net/http"
	"testing"

	"go-bootcamp/mock"
)

type fakeHandler struct {
	assertion func(*fakeHandler)
	Param1    int    `query:"param1"`
	Param2    string `query:"param2,required"`
}

func TestWrapHandler(t *testing.T) {
	fh := fakeHandler{
		assertion: func(h *fakeHandler) {
			if h.Param1 != 1 {
				t.Logf("Expected Param1 to be 1, got %v", h.Param1)
				t.Fail()
			}
			if h.Param2 != "test" {
				t.Logf("Expected Param2 to be [test], got %v", h.Param2)
				t.Fail()
			}
		},
	}

	wrapped := WrapHandler(&fh)
	response, _ := mock.NewFakeResponseWriter()
	request, _ := http.NewRequest("GET", "/test?param1=1&param2=test", mock.NewFakeReader(""))

	wrapped(response, request)
}

func TestRequiredParam(t *testing.T) {
	fh := fakeHandler{
		assertion: func(h *fakeHandler) {
			t.Logf("fakeHandler.ServeHTTP should not be reached with a missing required parameter")
			t.Fail()
		},
	}

	wrapped := WrapHandler(&fh)
	response, responseBody := mock.NewFakeResponseWriter()
	request, _ := http.NewRequest("GET", "/test?param1=1", mock.NewFakeReader(""))

	wrapped(response, request)
	expected := `{"errors":["Missing required field param2"]}`

	if !bytes.Equal([]byte(expected), responseBody()) {
		t.Logf("Response %s does not match expected %s", string(responseBody()), expected)
		t.Fail()
	}
}

func (i *fakeHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	i.assertion(i)
}
