package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleGetRootRR(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	handleGetRoot(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	tm := resp.Header.Get("x-datetime")
	if tm == "" {
		t.Error("Expected non-empty `x-datetime` header")
	}

	log.Printf("requested at: %s", tm)

	expected := "ok"
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Errorf("Expected body to be %s, but got %s", expected, string(b))
	}
}

func TestHandlePostRootRR(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key": "value"}`))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Error(err)
	}

	handleGetRoot(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, but got %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	defer resp.Body.Close()

	tm := resp.Header.Get("x-datetime")
	if tm == "" {
		t.Error("Expected non-empty `x-datetime` header")
	}

	log.Printf("requested at: %s", tm)
}

func TestHandleGetNotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/gophers/", nil)
	if err != nil {
		t.Error(err)
	}

	handleGetRoot(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, resp.StatusCode)
	}

	defer resp.Body.Close()
}

//
// disabled for now, check where it matters.
//
// func TestHandleGetRoot(t *testing.T) {
// 	server := httptest.NewServer(http.HandlerFunc(handleGetRoot))
// 	resp, err := http.Get(server.URL)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
// 	}

// 	defer resp.Body.Close()

// 	expected := "ok"
// 	b, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if string(b) != expected {
// 		t.Errorf("Expected body to be %s, but got %s", expected, string(b))
// 	}
// }
