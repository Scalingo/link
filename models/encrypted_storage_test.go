package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v2/config"
)

func Test_EncryptedStorageEndToEnd(t *testing.T) {
	ctx := context.Background()
	config := config.Config{
		SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
	}

	storage, err := NewEncryptedStorage(ctx, config)
	require.NoError(t, err)

	plainText := "Hello, World!"
	encryptedData, err := storage.Encrypt(ctx, plainText)
	require.NoError(t, err)

	var decryptedText string
	err = storage.Decrypt(ctx, encryptedData, &decryptedText)
	require.NoError(t, err)

	assert.Equal(t, plainText, decryptedText)
}

func Test_EncryptedStorageConstructor(t *testing.T) {
	ctx := context.Background()

	t.Run("if the key is too short", func(t *testing.T) {
		config := config.Config{
			SecretStorageEncryptionKey: "short-key",
		}

		storage, err := NewEncryptedStorage(ctx, config)
		require.Error(t, err)
		assert.Nil(t, storage)
	})

	t.Run("when everything is ok", func(t *testing.T) {
		config := config.Config{
			SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
		}

		storage, err := NewEncryptedStorage(ctx, config)
		require.NoError(t, err)
		require.NotNil(t, storage)
	})
}

func Test_EncryptedStorage_Decrypt(t *testing.T) {
	ctx := context.Background()
	config := config.Config{
		SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
	}

	storage, err := NewEncryptedStorage(ctx, config)
	require.NoError(t, err)

	t.Run("with an invalid nonce", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedData{
			Nonce: "invalid-nonce",
			Data:  "invalid-data",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode nonce")
	})

	t.Run("With an invalid data", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedData{
			Nonce: "000000000000000000000000",
			Data:  "invalid-data",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode cipher text")
	})

	t.Run("With a nonce of the wrong size", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedData{
			Nonce: "deadbeef",
			Data:  "deadbeef",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "nonce size is incorrect")
	})

	t.Run("With a cipher text that is invalid", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedData{
			Nonce: "000000000000000000000000",
			Data:  "deadbeef",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "message authentication failed")
	})

	t.Run("With a valid cipher text", func(t *testing.T) {
		encryptedData, err := storage.Encrypt(ctx, "Hello, World!")
		require.NoError(t, err)

		var decryptedText string
		err = storage.Decrypt(ctx, encryptedData, &decryptedText)
		require.NoError(t, err)
		assert.Equal(t, "Hello, World!", decryptedText)
	})
}
