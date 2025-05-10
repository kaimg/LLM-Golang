package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

var secretKey []byte

func init() {
	// Get the encryption key from environment variable or generate one
	encodedKey := os.Getenv("ENCRYPTION_KEY")
	if encodedKey == "" {
		// Generate a random 32-byte key for AES-256
		secretKey = make([]byte, 32)
		if _, err := io.ReadFull(rand.Reader, secretKey); err != nil {
			panic(fmt.Errorf("failed to generate encryption key: %v", err))
		}
	} else {
		// Decode the base64-encoded key
		var err error
		secretKey, err = base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			panic(fmt.Errorf("failed to decode base64 encryption key: %v", err))
		}
		if len(secretKey) != 32 {
			panic(fmt.Sprintf("ENCRYPTION_KEY must decode to exactly 32 bytes, got %d bytes", len(secretKey)))
		}
	}
}

// EncryptAPIKey encrypts the API key using AES-256-GCM
func EncryptAPIKey(apiKey string) (string, error) {
	if apiKey == "" {
		return "", nil
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAPIKey decrypts the encrypted API key
func DecryptAPIKey(encryptedKey string) (string, error) {
	if encryptedKey == "" {
		return "", nil
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plaintext), nil
} 