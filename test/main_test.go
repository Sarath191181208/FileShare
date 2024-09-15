package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"sarath/backend_project/cmd/api"
	"sarath/backend_project/internal/cache"
	"sarath/backend_project/internal/data"
	filestore "sarath/backend_project/internal/file_store"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Mock structures
type MockLogger struct{}

func (l *MockLogger) Printf(format string, args ...interface{}) {
  fmt.Printf(format, args...)
}

// Setup the test database
func setupTestDB() (*sql.DB, error) {
	// Create an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	// Create necessary tables and insert test data here
	// e.g. db.Exec(`CREATE TABLE users (...)`)
	db.Exec(`
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  email TEXT UNIQUE NOT NULL,
  password_hash BLOB NOT NULL,
  version INTEGER NOT NULL DEFAULT 1
);`)

	db.Exec(`
CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    upload_date TEXT NOT NULL DEFAULT (datetime('now')),
    size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    file_url TEXT NOT NULL
);
    `)

	return db, nil
}


func TestRoutes(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("could not set up test database: %v", err)
	}
	defer db.Close()

	app := &api.Application{
    Config: api.Config{
      Jwt: struct{Secret string} {
        Secret: "TH",
      },
    },
		Logger:    log.New(os.Stdout, "", log.Ldate|log.Ltime),
		Models:    data.NewModels(db), // Replace with an actual implementation using db
		FileStore: &filestore.MockFileStore{},
		Cache:     &cache.MockCache{},
	}

	router := app.Routes()

	t.Run("POST /register", func(t *testing.T) {
		userData := map[string]string{"email": "testuser@gmail.com", "password": "testpass"}
		jsonData, _ := json.Marshal(userData)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("POST /login", func(t *testing.T) {
		loginData := map[string]string{"email": "testuser@gmail.com", "password": "testpass"}
		jsonData, _ := json.Marshal(loginData)

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
