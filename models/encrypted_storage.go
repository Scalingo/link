package models

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v2/config"
)

type EncryptedData struct {
	Nonce string `json:"nonce"`
	Data  string `json:"data"`
}

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

	gcmInstance, err := s.initializeAESGCMCipher(ctx)
	if err != nil {
		return EncryptedData{}, errors.Wrap(ctx, err, "initialize AES GCM cipher")
	}

	// Generate a random nonce
	nonce := make([]byte, gcmInstance.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return EncryptedData{}, errors.Wrap(ctx, err, "generate nonce")
	}

	result := gcmInstance.Seal(nil, nonce, dataBytes, nil)
	return EncryptedData{
		Nonce: hex.EncodeToString(nonce),
		Data:  hex.EncodeToString(result),
	}, nil
}

func (s *encryptedStorage) Decrypt(ctx context.Context, data EncryptedData, v any) error {
	gcmInstance, err := s.initializeAESGCMCipher(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "initialize AES GCM cipher")
	}

	nonce, err := hex.DecodeString(data.Nonce)
	if err != nil {
		return errors.Wrap(ctx, err, "decode nonce")
	}

	cipherText, err := hex.DecodeString(data.Data)
	if err != nil {
		return errors.Wrap(ctx, err, "decode cipher text")
	}

	if len(nonce) != gcmInstance.NonceSize() {
		return errors.New(ctx, "nonce size is incorrect")
	}

	plaintext, err := gcmInstance.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return errors.Wrap(ctx, err, "decrypt data")
	}

	err = json.Unmarshal(plaintext, v)
	if err != nil {
		return errors.Wrap(ctx, err, "unmarshal decrypted data")
	}

	return nil
}

func (s *encryptedStorage) initializeAESGCMCipher(ctx context.Context) (cipher.AEAD, error) {
	if len(s.secretKey) != 32 {
		return nil, errors.New(ctx, "secret key must be 32 bytes long")
	}

	block, err := aes.NewCipher(s.secretKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create AES cipher")
	}
	gcmInstance, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create GCM instance")
	}

	return gcmInstance, nil
}
