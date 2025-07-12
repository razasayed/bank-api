package tests

import (
	"bank-api/db"
	"bank-api/handlers"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	// Resolve absolute path to project root
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(filename), "..")
	envPath := filepath.Join(projectRoot, ".env.test")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Failed to load .env.test from %s: %v", envPath, err)
	}

	db.Connect()
	code := m.Run()
	_ = db.DB.Close()
	os.Exit(code)
}

func ResetDB() {
	_, err := db.DB.Exec(`
		TRUNCATE accounts RESTART IDENTITY CASCADE;
		TRUNCATE transactions RESTART IDENTITY CASCADE;
	`)
	if err != nil {
		log.Fatalf("failed to reset test DB: %v", err)
	}
}

func seedAccount(t *testing.T, documentNumber string) (int, string) {
	body := []byte(`{"document_number":"` + documentNumber + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handlers.CreateAccountHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to seed account: got status %d", w.Code)
	}

	var resp struct {
		AccountID      int    `json:"account_id"`
		DocumentNumber string `json:"document_number"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("seedAccount: Invalid JSON response: %v", err)
	}
	if resp.DocumentNumber != documentNumber {
		t.Fatalf("seedAccount: expected document number %q, got %q", documentNumber, resp.DocumentNumber)
	}
	return resp.AccountID, resp.DocumentNumber
}
