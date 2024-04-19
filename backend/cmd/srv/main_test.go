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
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/zyghq/zyg/internal/model"
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

	t.Logf("requested at: %v", tm)

	expected := "ok"
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Fatalf("Expected body to be %s, but got %s", expected, string(b))
	}
}

func DB(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	var db *pgxpool.Pool

	pgConnStr, pgConnStatus := os.LookupEnv("POSTGRES_URI")
	if !pgConnStatus {
		err = fmt.Errorf("env `POSTGRES_URI` is not set")
		return db, err
	}

	db, err = pgxpool.New(ctx, pgConnStr)
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
		t.Fatalf("could not connect to database %v", err)
	}

	defer db.Close()

	handler := handleGetWorkspaces(ctx, db)
	handler.ServeHTTP(rr, req)

	resp := rr.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	defer resp.Body.Close()

	var workspaces []model.Workspace

	err = json.NewDecoder(resp.Body).Decode(&workspaces)

	if err != nil {
		t.Fatalf("could not decode json response %v", err)
	}
}

func TestHeadTime(t *testing.T) {
	resp, err := http.Head("https://www.time.gov/")

	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	now := time.Now().Round(time.Second)
	date := resp.Header.Get("Date")

	if date == "" {
		t.Fatal("Expected non-empty `Date` header")
	}

	dt, err := time.Parse(time.RFC1123, date)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("now: %v, date: %v", now, dt)
}
