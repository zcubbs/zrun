// Package vault provides a set of functions to interact with the vault.
/*
Copyright © 2023 zcubbs https://github.com/zcubbs
*/
package vault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// Secret struct will hold the data for our secret
type Secret struct {
	Key   string
	Value string
}

// SecretVault struct will hold our secrets in memory
type SecretVault struct {
	secrets map[string]Secret
	mutex   sync.RWMutex
}

// AddSecret is a function to add a secret to our SecretVault
func (s *SecretVault) AddSecret(key, value string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.secrets[key] = Secret{
		Key:   key,
		Value: value,
	}

	return nil
}

// GetSecret is a function to get a secret from our SecretVault
func (s *SecretVault) GetSecret(key string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	secret, ok := s.secrets[key]
	if !ok {
		return "", errors.New("no secret found with that key")
	}

	return secret.Value, nil
}

// NewSecretVault reads data from the file and decrypts it if exists
func NewSecretVault() (*SecretVault, error) {
	v := &SecretVault{
		secrets: make(map[string]Secret),
	}

	filename := os.Getenv("VAULT_FILE")
	key := os.Getenv("VAULT_KEY")
	if filename == "" || key == "" {
		return nil, errors.New("missing VAULT_FILE or VAULT_KEY env variable")
	}

	encrypted, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If the file does not exist, return the new empty vault
			return v, nil
		}
		return nil, fmt.Errorf("error reading vault file: %w", err)
	}

	data, err := decrypt(encrypted, []byte(key))
	if err != nil {
		return nil,
			fmt.Errorf("error decrypting vault: %w", err)
	}

	err = gob.NewDecoder(bytes.NewReader(data)).Decode(&v.secrets)
	if err != nil {
		return nil,
			fmt.Errorf("error decoding vault: %w", err)
	}

	return v, nil
}

func InitializeVaultWithRandomKey(filename string) (*SecretVault, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil,
			fmt.Errorf("error reading key from random source: %w", err)
	}

	// convert key to hex string
	//keyStr := hex.EncodeToString(key)

	// set environment variables
	err = os.Setenv("VAULT_FILE", filename)
	if err != nil {
		return nil, err
	}
	err = os.Setenv("VAULT_KEY", string(key))
	if err != nil {
		return nil, err
	}

	// return new SecretVault
	return NewSecretVault()
}

func (s *SecretVault) Save() error {
	filename := os.Getenv("VAULT_FILE")
	key := os.Getenv("VAULT_KEY")

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(s.secrets)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(buf.Bytes(), []byte(key))
	if err != nil {
		return err
	}

	return os.WriteFile(filename, encrypted, 0600)
}

// Helper functions for encryption and decryption
func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
