package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMainPackageCompiles(t *testing.T) {
	t.Log("cmd package compiles and can be imported")
}

func TestParseServerConfigRejectsAuthFlagsWithoutHTTP(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "keys.db")
	tests := []struct {
		name string
		args []string
	}{
		{name: "auth", args: []string{"-auth", databasePath}},
		{name: "rate limit", args: []string{"-rate-limit", "1"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := parseServerConfig(test.args); err == nil {
				t.Fatal("expected flag combination to be rejected")
			}
			if _, err := os.Stat(databasePath); !os.IsNotExist(err) {
				t.Fatalf("auth database was opened during flag parsing: %v", err)
			}
		})
	}
}

func TestParseServerConfigValidatesAuthFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "whitespace auth path", args: []string{"-http", "127.0.0.1:8080", "-auth", "  "}},
		{name: "negative rate limit", args: []string{"-http", "127.0.0.1:8080", "-rate-limit", "-1"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := parseServerConfig(test.args); err == nil {
				t.Fatal("expected invalid flag value to be rejected")
			}
		})
	}
}

func TestParseServerConfigAcceptsHTTPConfiguration(t *testing.T) {
	config, err := parseServerConfig([]string{"-http", "127.0.0.1:8080", "-auth", "/tmp/keys.db", "-rate-limit", "5"})
	if err != nil {
		t.Fatalf("parse server config: %v", err)
	}
	if config.httpAddr != "127.0.0.1:8080" {
		t.Errorf("http address: got %q, want %q", config.httpAddr, "127.0.0.1:8080")
	}
	if config.authPath != "/tmp/keys.db" {
		t.Errorf("auth path: got %q, want %q", config.authPath, "/tmp/keys.db")
	}
	if config.rateLimit != 5 {
		t.Errorf("rate limit: got %d, want %d", config.rateLimit, 5)
	}
}
