package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetRootRR(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}

	handleGetRoot(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	expected := "ok"
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Errorf("Expected body to be %s, but got %s", expected, string(b))
	}
}

func TestHandleGetRoot(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleGetRoot))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	expected := "ok"
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Errorf("Expected body to be %s, but got %s", expected, string(b))
	}
}
