package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestHandleGetIndexRR(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	handleGetIndex(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	tm := resp.Header.Get("x-datetime")
	if tm == "" {
		t.Fatal("Expected non-empty `x-datetime` header")
	}

	t.Log("requested at: ", tm)

	expected := "ok"
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Fatalf("Expected body to be %s, but got %s", expected, string(b))
	}
}

func DB(ctx context.Context) (*pgx.Conn, error) {
	var err error
	var db *pgx.Conn

	pgConnStr, pgConnStatus := os.LookupEnv("POSTGRES_URI")
	if !pgConnStatus {
		err = fmt.Errorf("env `POSTGRES_URI` is not set")
		return db, err
	}

	db, err = pgx.Connect(ctx, pgConnStr)
	if err != nil {
		err = fmt.Errorf("failed to connect database: %v", err)
		return db, err
	}
	return db, err
}

func TestHandleGetWorkspacesRR(t *testing.T) {
	ctx := context.Background()
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/workspaces/", nil)
	if err != nil {
		t.Fatal(err)
	}

	db, err := DB(ctx)
	if err != nil {
		t.Fatalf("could not make connection to DB %v", err)
	}

	handler := handleGetWorkspaces(ctx, db)
	handler.ServeHTTP(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	var workspaces []Workspace

	err = json.NewDecoder(resp.Body).Decode(&workspaces)

	if err != nil {
		t.Fatalf("could not decode response %v", err)
	}
}
