package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/config"
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

	t.Run("with an invalid type", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedDataLink{
			Type: "invalid-type",
			Data: "invalid-data",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported encryption type")
	})

	t.Run("With an invalid data", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedDataLink{
			Type: EncryptedDataTypeAESCFB,
			Data: "invalid-data",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "decode cipher text")
	})

	t.Run("With a cipher text that is invalid", func(t *testing.T) {
		err := storage.Decrypt(ctx, EncryptedDataLink{
			Type: EncryptedDataTypeAESCFB,
			Data: "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
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
