package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 16
	keySize    = 32
	iterations = 100000
)

// Encrypt encrypts the plaintext using AES-GCM with a key derived from the password via PBKDF2
func Encrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Final blob: [salt][nonce][ciphertext]
	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt decrypts the ciphertext using a key derived from the password
func Decrypt(data []byte, password string) ([]byte, error) {
	if len(data) < saltSize+12 { // 12 is the default GCM nonce size
		return nil, errors.New("data too short")
	}

	salt := data[:saltSize]
	nonce := data[saltSize : saltSize+12]
	ciphertext := data[saltSize+12:]

	key := pbkdf2.Key([]byte(password), salt, iterations, keySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}
