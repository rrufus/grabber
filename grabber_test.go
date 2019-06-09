package grabber

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGrabOneUrl(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(`{"testResult": "hi"}`))
	}))
	defer testServer.Close()

	search := &Item{Key: "test", Url: testServer.URL}

	result, err := Grab(testServer.Client(), search)
	if err != nil {
		t.Fatal(err)
	}

	if string(result["test"]) != `{"testResult": "hi"}` {
		t.Fatalf(`Expected result to by {"testResult": "hi"} but was [%v]`, string(result["text"]))
	}
}
func TestGrabTwoUrls(t *testing.T) {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte(`{"testResult": "hi"}`))
	}))
	defer testServer.Close()

	search1 := &Item{Key: "test1", Url: testServer.URL}
	search2 := &Item{Key: "test2", Url: testServer.URL}

	result, err := Grab(testServer.Client(), search1, search2)
	if err != nil {
		t.Fatal(err)
	}

	if string(result["test1"]) != `{"testResult": "hi"}` {
		t.Fatalf(`Expected result to by {"testResult": "hi"} but was [%v]`, string(result["test1"]))
	}

	if string(result["test2"]) != `{"testResult": "hi"}` {
		t.Fatalf(`Expected result to by {"testResult": "hi"} but was [%v]`, string(result["test2"]))
	}
}

func TestResponseError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusBadRequest)
	}))
	defer testServer.Close()

	search := &Item{Key: "test", Url: testServer.URL}

	_, err := Grab(testServer.Client(), search)
	if err == nil {
		t.Fatal("Expected an error got none")
	}
}
