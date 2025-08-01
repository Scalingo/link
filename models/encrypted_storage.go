package models

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"

	"github.com/Scalingo/go-utils/crypto"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/config"
)

type EncryptedDataLink struct {
	ID         string `json:"id"`
	EndpointID string `json:"endpoint_id"`

	// Depracated: This is here to keep compatibility with old storage method
	Type string `json:"type,omitempty"`
	Data string `json:"data,omitempty"`
	Hash string `json:"hash,omitempty"`
}

type EncryptedData struct {
	ID         string `json:"id"`
	EndpointID string `json:"endpoint_id"`
	Type       string `json:"type"`
	Data       string `json:"data"`
	Hash       string `json:"hash"`
}

const (
	// Deprecated: use EncryptedDataTypeAESCFBSha512 instead
	EncryptedDataTypeAESCFB = "aes-cfb"

	EncryptedDataTypeAESCFBSha512 = "aes-cfb-sha512"
)

type EncryptedStorage interface {
	Encrypt(ctx context.Context, endpointID string, data any) (EncryptedDataLink, error)
	Decrypt(ctx context.Context, data EncryptedDataLink, v any) error
	Cleanup(ctx context.Context, endpointID string) error
}

// Implements a GCM AES-256 encryption/decryption
type encryptedStorage struct {
	secretKey     []byte   // AES-256 key
	alternateKeys [][]byte // Optional alternate keys for decryption
	storage       Storage
}

func NewEncryptedStorage(ctx context.Context, config config.Config, storage Storage) (EncryptedStorage, error) {
	if len(config.SecretStorageEncryptionKey) < 32 {
		return nil, errors.New(ctx, "SecretStorageEncryptionKey must be at least 32 characters long")
	}

	key := sha256.Sum256([]byte(config.SecretStorageEncryptionKey))
	if len(key) != 32 {
		return nil, errors.New(ctx, "SecretStorageEncryptionKey must be 32 bytes long after hashing")
	}

	altKeys := make([][]byte, 0, len(config.SecretStorageAlternateKeys))

	for _, altKey := range config.SecretStorageAlternateKeys {
		if len(altKey) < 32 {
			return nil, errors.New(ctx, "each SecretStorageAlternateKey must be at least 32 characters long")
		}
		altKeyBytes := sha256.Sum256([]byte(altKey))
		if len(altKeyBytes) != 32 {
			return nil, errors.New(ctx, "each SecretStorageAlternateKey must be 32 bytes long after hashing")
		}
		altKeys = append(altKeys, altKeyBytes[:])
	}

	return &encryptedStorage{
		secretKey:     key[:], // Convert [32]byte to []byte
		alternateKeys: altKeys,
		storage:       storage,
	}, nil
}

func (s *encryptedStorage) Encrypt(ctx context.Context, endpointId string, data any) (EncryptedDataLink, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return EncryptedDataLink{}, errors.Wrap(ctx, err, "marshal data to JSON")
	}

	result, err := crypto.Encrypt(s.secretKey, dataBytes)
	if err != nil {
		return EncryptedDataLink{}, errors.Wrap(ctx, err, "encrypt data")
	}

	hash := sha512.Sum512(dataBytes)

	encryptedData := EncryptedData{
		Data: hex.EncodeToString(result),
		Type: EncryptedDataTypeAESCFBSha512,
		Hash: hex.EncodeToString(hash[:]),
	}

	link, err := s.storage.UpsertEncryptedData(ctx, endpointId, encryptedData)
	if err != nil {
		return EncryptedDataLink{}, errors.Wrap(ctx, err, "add encrypted data to storage")
	}

	return link, nil
}

func (s *encryptedStorage) Decrypt(ctx context.Context, link EncryptedDataLink, v any) error {
	var data EncryptedData
	var err error
	if link.ID == "" {
		// All encrypted data should be stored in the storage
		// However some old versions (LinK v3.0.0 to v3.0.2) stored encrypted data in the endpoint directly.
		// We keep this code for backward compatibility.
		data = EncryptedData{
			Type: link.Type,
			Data: link.Data,
			Hash: link.Hash,
		}
	} else {
		data, err = s.storage.GetEncryptedData(ctx, link.EndpointID, link.ID)
		if err != nil {
			return errors.Wrap(ctx, err, "get encrypted data from storage")
		}
	}

	if data.Type != EncryptedDataTypeAESCFB && data.Type != EncryptedDataTypeAESCFBSha512 {
		return errors.New(ctx, "unsupported encryption type: "+data.Type)
	}
	var plaintext []byte

	keys := append([][]byte{s.secretKey}, s.alternateKeys...)

	for _, key := range keys {
		plaintext, err = s.decryptWithKey(ctx, data, key)
		if err == nil {
			break
		}
	}
	if err != nil {
		return errors.Wrap(ctx, err, "decrypt data with all keys")
	}

	if data.Type == EncryptedDataTypeAESCFBSha512 {
		hash := sha512.Sum512(plaintext)
		if hex.EncodeToString(hash[:]) != data.Hash {
			return errors.New(ctx, "hash mismatch after decryption")
		}
	}

	err = json.Unmarshal(plaintext, v)
	if err != nil {
		return errors.Wrap(ctx, err, "unmarshal decrypted data")
	}

	return nil
}

func (s *encryptedStorage) decryptWithKey(ctx context.Context, data EncryptedData, key []byte) ([]byte, error) {
	cipherText, err := hex.DecodeString(data.Data)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "decode cipher text")
	}

	plaintext, err := crypto.Decrypt(key, cipherText)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "decrypt data with alternate key")
	}

	if data.Type == EncryptedDataTypeAESCFBSha512 {
		hash := sha512.Sum512(plaintext)
		if hex.EncodeToString(hash[:]) != data.Hash {
			return nil, errors.New(ctx, "hash mismatch after decryption")
		}
	}

	return plaintext, nil
}

func (s *encryptedStorage) Cleanup(ctx context.Context, endpointID string) error {
	err := s.storage.RemoveEncryptedDataForEndpoint(ctx, endpointID)
	if err != nil {
		return errors.Wrap(ctx, err, "remove encrypted data for endpoint")
	}

	return nil
}
