package auth

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStorePersistsAPIKeyAcrossRestart(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "keys.db")
	store, err := OpenDatabase(databasePath)
	if err != nil {
		t.Fatalf("open auth database: %v", err)
	}
	apiKey, err := store.Register(" User@Example.COM ")
	if err != nil {
		t.Fatalf("register API key: %v", err)
	}
	if err := store.Close(); err != nil {
		t.Fatalf("close auth database: %v", err)
	}

	databaseBytes, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("read auth database: %v", err)
	}
	if strings.Contains(string(databaseBytes), apiKey) {
		t.Error("auth database contains the plaintext API key")
	}

	store, err = OpenDatabase(databasePath)
	if err != nil {
		t.Fatalf("reopen auth database: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("close auth database: %v", err)
		}
	}()

	email, ok, err := store.Authenticate(apiKey)
	if err != nil {
		t.Fatalf("authenticate API key: %v", err)
	}
	if !ok {
		t.Fatal("persisted API key was not found")
	}
	if email != "user@example.com" {
		t.Errorf("authenticated email: got %q, want %q", email, "user@example.com")
	}

	_, err = store.Register("user@example.com")
	if !errors.Is(err, ErrEmailExists) {
		t.Fatalf("duplicate register error: got %v, want %v", err, ErrEmailExists)
	}
}

func TestStoreRejectsInvalidAPIKeys(t *testing.T) {
	store, err := OpenDatabase(filepath.Join(t.TempDir(), "keys.db"))
	if err != nil {
		t.Fatalf("open auth database: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			t.Errorf("close auth database: %v", err)
		}
	}()

	for _, apiKey := range []string{"", "not-hex", strings.Repeat("a", 63), strings.Repeat("a", 65)} {
		_, ok, err := store.Authenticate(apiKey)
		if err != nil {
			t.Fatalf("authenticate invalid API key: %v", err)
		}
		if ok {
			t.Errorf("invalid API key %q was accepted", apiKey)
		}
	}
}
