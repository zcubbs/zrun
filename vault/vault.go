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
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

// Secret struct will hold the data for our secret
type Secret struct {
	key  string
	hash []byte
}

// SecretVault struct will hold our secrets in memory
type SecretVault struct {
	secrets map[string]Secret
	mutex   sync.RWMutex
}

// AddSecret is a function to add a secret to our SecretVault
func (s *SecretVault) AddSecret(key, password string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.secrets[key] = Secret{
		key:  key,
		hash: hash,
	}

	return nil
}

// GetSecret is a function to get a secret from our SecretVault
func (s *SecretVault) GetSecret(key, password string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	secret, ok := s.secrets[key]
	if !ok {
		return "", errors.New("no secret found with that key")
	}

	err := bcrypt.CompareHashAndPassword(secret.hash, []byte(password))
	if err != nil {
		return "", errors.New("invalid password for secret")
	}

	return secret.key, nil
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

	encrypted, err := ioutil.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// If the file does not exist, return the new empty vault
			return v, nil
		}
		return nil, err
	}

	data, err := decrypt(encrypted, []byte(key))
	if err != nil {
		return nil, err
	}

	err = gob.NewDecoder(bytes.NewReader(data)).Decode(&v.secrets)
	if err != nil {
		return nil, err
	}

	return v, nil
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

	return ioutil.WriteFile(filename, encrypted, 0600)
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
