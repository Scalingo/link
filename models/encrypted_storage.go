package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/Scalingo/go-utils/crypto"
	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/config"
)

type EncryptedData struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

const (
	EncryptedDataTypeAESCFB = "aes-cfb"
)

type EncryptedStorage interface {
	Encrypt(ctx context.Context, data any) (EncryptedData, error)
	Decrypt(ctx context.Context, data EncryptedData, v any) error
}

// Implements a GCM AES-256 encryption/decryption
type encryptedStorage struct {
	secretKey []byte // AES-256 key
}

func NewEncryptedStorage(ctx context.Context, config config.Config) (EncryptedStorage, error) {
	if len(config.SecretStorageEncryptionKey) < 32 {
		return nil, errors.New(ctx, "SecretStorageEncryptionKey must be at least 32 characters long")
	}

	key := sha256.Sum256([]byte(config.SecretStorageEncryptionKey))
	if len(key) != 32 {
		return nil, errors.New(ctx, "SecretStorageEncryptionKey must be 32 bytes long after hashing")
	}

	return &encryptedStorage{
		secretKey: key[:], // Convert [32]byte to []byte
	}, nil
}

func (s *encryptedStorage) Encrypt(ctx context.Context, data any) (EncryptedData, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return EncryptedData{}, errors.Wrap(ctx, err, "marshal data to JSON")
	}

	result, err := crypto.Encrypt(s.secretKey, dataBytes)
	if err != nil {
		return EncryptedData{}, errors.Wrap(ctx, err, "encrypt data")
	}

	return EncryptedData{
		Data: hex.EncodeToString(result),
		Type: EncryptedDataTypeAESCFB,
	}, nil
}

func (s *encryptedStorage) Decrypt(ctx context.Context, data EncryptedData, v any) error {
	if data.Type != EncryptedDataTypeAESCFB {
		return errors.New(ctx, "unsupported encryption type: "+data.Type)
	}

	cipherText, err := hex.DecodeString(data.Data)
	if err != nil {
		return errors.Wrap(ctx, err, "decode cipher text")
	}

	plaintext, err := crypto.Decrypt(s.secretKey, cipherText)
	if err != nil {
		return errors.Wrap(ctx, err, "decrypt data")
	}

	err = json.Unmarshal(plaintext, v)
	if err != nil {
		return errors.Wrap(ctx, err, "unmarshal decrypted data")
	}

	return nil
}
