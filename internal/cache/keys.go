package cache

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LicenseKey represents a stored license key
type LicenseKey struct {
	ID      string    `json:"id"`
	Label   string    `json:"label"`
	Key     string    `json:"key"`
	Created time.Time `json:"created"`
}

// KeyVault manages encrypted license keys
type KeyVault struct {
	filePath string
	keys     []LicenseKey
}

// NewKeyVault creates a new key vault
func NewKeyVault(filePath string) (*KeyVault, error) {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create vault directory: %w", err)
	}

	kv := &KeyVault{
		filePath: filePath,
	}

	// Load or create vault
	if err := kv.load(); err != nil {
		return nil, err
	}

	return kv, nil
}

// Add adds a new license key
func (kv *KeyVault) Add(label, key string) (string, error) {
	// Validate key format (should start with cfxk_)
	if len(key) < 10 || key[:5] != "cfxk_" {
		return "", fmt.Errorf("invalid license key format")
	}

	// Check if key already exists
	for _, existingKey := range kv.keys {
		if existingKey.Key == key {
			return "", fmt.Errorf("key already exists")
		}
	}

	// Create new key entry
	id := uuid.New().String()
	licenseKey := LicenseKey{
		ID:      id,
		Label:   label,
		Key:     key,
		Created: time.Now(),
	}

	kv.keys = append(kv.keys, licenseKey)

	if err := kv.save(); err != nil {
		return "", err
	}

	return id, nil
}

// Remove removes a license key by ID
func (kv *KeyVault) Remove(id string) error {
	for i, key := range kv.keys {
		if key.ID == id {
			kv.keys = append(kv.keys[:i], kv.keys[i+1:]...)
			return kv.save()
		}
	}

	return fmt.Errorf("key not found")
}

// Get retrieves a license key by ID
func (kv *KeyVault) Get(id string) (*LicenseKey, error) {
	for i, key := range kv.keys {
		if key.ID == id {
			return &kv.keys[i], nil
		}
	}

	return nil, fmt.Errorf("key not found")
}

// List returns all license keys (with masked keys for display)
func (kv *KeyVault) List() []LicenseKey {
	return kv.keys
}

// Count returns the number of stored keys
func (kv *KeyVault) Count() int {
	return len(kv.keys)
}

// MaskKey returns a masked version of a key for display
func MaskKey(key string) string {
	if len(key) < 15 {
		return "****"
	}

	// Show first 5 chars (cfxk_) and last 4 chars
	return key[:5] + strings.Repeat("*", len(key)-9) + key[len(key)-4:]
}

// load loads the vault from disk (encrypted)
func (kv *KeyVault) load() error {
	// If vault doesn't exist, create empty
	if _, err := os.Stat(kv.filePath); os.IsNotExist(err) {
		kv.keys = []LicenseKey{}
		return kv.save()
	}

	// Read encrypted data
	encrypted, err := os.ReadFile(kv.filePath)
	if err != nil {
		return fmt.Errorf("failed to read vault: %w", err)
	}

	// Decrypt
	data, err := kv.decrypt(encrypted)
	if err != nil {
		return fmt.Errorf("failed to decrypt vault: %w", err)
	}

	// Parse JSON
	var keys []LicenseKey
	if err := json.Unmarshal(data, &keys); err != nil {
		return fmt.Errorf("failed to parse vault: %w", err)
	}

	kv.keys = keys
	return nil
}

// save saves the vault to disk (encrypted)
func (kv *KeyVault) save() error {
	// Marshal to JSON
	data, err := json.MarshalIndent(kv.keys, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	// Encrypt
	encrypted, err := kv.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}

	// Write to file
	if err := os.WriteFile(kv.filePath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write vault: %w", err)
	}

	return nil
}

// encrypt encrypts data using AES-256-GCM
func (kv *KeyVault) encrypt(plaintext []byte) ([]byte, error) {
	key := kv.getMachineKey()

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

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-256-GCM
func (kv *KeyVault) decrypt(ciphertext []byte) ([]byte, error) {
	key := kv.getMachineKey()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// getMachineKey derives a machine-specific encryption key
func (kv *KeyVault) getMachineKey() []byte {
	// Get machine ID (hostname for simplicity)
	hostname, _ := os.Hostname()

	// Use vault file path as additional entropy
	combined := hostname + kv.filePath

	// SHA-256 hash for 32-byte key
	hash := sha256.Sum256([]byte(combined))
	return hash[:]
}
