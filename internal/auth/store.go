package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	bolt "go.etcd.io/bbolt"
)

var (
	emailBucketName = []byte("emails")
	keyBucketName   = []byte("keys")
)

// ErrEmailExists indicates that registration was attempted for an email that is
// already registered.
var ErrEmailExists = errors.New("email already exists")

// Store persists normalized email addresses and hashes of API keys in BoltDB.
type Store struct {
	db *bolt.DB
}

// OpenDatabase opens or creates a Bolt database for API key authentication.
func OpenDatabase(path string) (*Store, error) {
	db, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open auth database: %w", err)
	}

	return &Store{db: db}, nil
}

// Register creates an API key for an email address and returns the plaintext
// key. The plaintext key is not persisted.
func (s *Store) Register(email string) (string, error) {
	normalizedEmail := normalizeEmail(email)
	if normalizedEmail == "" {
		return "", errors.New("email is empty")
	}

	rawKey := make([]byte, 32)
	if _, err := rand.Read(rawKey); err != nil {
		return "", fmt.Errorf("generate API key: %w", err)
	}
	keyHash := hashKey(rawKey)

	err := s.db.Update(func(tx *bolt.Tx) error {
		emails, err := tx.CreateBucketIfNotExists(emailBucketName)
		if err != nil {
			return fmt.Errorf("create emails bucket: %w", err)
		}
		keys, err := tx.CreateBucketIfNotExists(keyBucketName)
		if err != nil {
			return fmt.Errorf("create keys bucket: %w", err)
		}
		if emails.Get([]byte(normalizedEmail)) != nil {
			return ErrEmailExists
		}
		if err := emails.Put([]byte(normalizedEmail), []byte(keyHash)); err != nil {
			return fmt.Errorf("store email: %w", err)
		}
		if err := keys.Put([]byte(keyHash), []byte(normalizedEmail)); err != nil {
			return fmt.Errorf("store API key hash: %w", err)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("register API key: %w", err)
	}

	return hex.EncodeToString(rawKey), nil
}

// Authenticate returns the normalized email associated with an API key.
func (s *Store) Authenticate(apiKey string) (string, bool, error) {
	rawKey, err := hex.DecodeString(apiKey)
	if err != nil || len(rawKey) != 32 {
		return "", false, nil
	}
	keyHash := hashKey(rawKey)

	var email string
	err = s.db.View(func(tx *bolt.Tx) error {
		keys := tx.Bucket(keyBucketName)
		if keys == nil {
			return nil
		}
		value := keys.Get([]byte(keyHash))
		if value == nil {
			return nil
		}
		email = string(value)
		return nil
	})
	if err != nil {
		return "", false, fmt.Errorf("authenticate API key: %w", err)
	}
	if email == "" {
		return "", false, nil
	}
	return email, true, nil
}

// Close closes the underlying database.
func (s *Store) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close auth database: %w", err)
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func hashKey(rawKey []byte) string {
	hash := sha256.Sum256(rawKey)
	return hex.EncodeToString(hash[:])
}
