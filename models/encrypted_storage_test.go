package models

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-utils/crypto"
	"github.com/Scalingo/link/v3/config"
)

func Test_EncryptedStorageEndToEnd(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctx := context.Background()
	storage := NewMockStorage(ctrl)
	config := config.Config{
		SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
	}

	encryptedStorage, err := NewEncryptedStorage(ctx, config, storage)
	require.NoError(t, err)

	var storedData EncryptedData
	storage.EXPECT().UpsertEncryptedData(gomock.Any(), "endpoint-id", gomock.Any()).DoAndReturn(func(ctx context.Context, endpointID string, data EncryptedData) (EncryptedDataLink, error) {
		storedData = data
		return EncryptedDataLink{
			ID:         "data-id",
			EndpointID: endpointID,
		}, nil
	})

	storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").DoAndReturn(func(ctx context.Context, _, id string) (EncryptedData, error) {
		return storedData, nil
	})

	plainText := "Hello, World!"
	encryptedData, err := encryptedStorage.Encrypt(ctx, "endpoint-id", plainText)
	require.NoError(t, err)

	var decryptedText string
	err = encryptedStorage.Decrypt(ctx, encryptedData, &decryptedText)
	require.NoError(t, err)

	assert.Equal(t, plainText, decryptedText)
}

func Test_EncryptedStorageConstructor(t *testing.T) {
	ctx := context.Background()

	t.Run("if the key is too short", func(t *testing.T) {
		config := config.Config{
			SecretStorageEncryptionKey: "short-key",
		}

		storage, err := NewEncryptedStorage(ctx, config, nil)
		require.Error(t, err)
		assert.Nil(t, storage)
	})

	t.Run("when everything is ok", func(t *testing.T) {
		config := config.Config{
			SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
		}

		storage, err := NewEncryptedStorage(ctx, config, nil)
		require.NoError(t, err)
		require.NotNil(t, storage)
	})
}

func Test_EncryptedStorage_Decrypt(t *testing.T) {

	ctx := t.Context()

	cfg := config.Config{
		SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
	}

	t.Run("with an invalid type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		storage := NewMockStorage(ctrl)

		storedData := EncryptedData{
			Type: "invalid-type",
			Data: "invalid-data",
		}
		storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

		encryptedStorage, err := NewEncryptedStorage(ctx, cfg, storage)
		require.NoError(t, err)

		err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
			ID:         "data-id",
			EndpointID: "endpoint-id",
		}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported encryption type")
	})

	t.Run("AES CFB decryption", func(t *testing.T) {
		t.Run("With an invalid data", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := NewMockStorage(ctrl)

			storedData := EncryptedData{
				Type: EncryptedDataTypeAESCFB,
				Data: "invalid-data",
			}
			storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

			encryptedStorage, err := NewEncryptedStorage(ctx, cfg, storage)
			require.NoError(t, err)

			err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
				ID:         "data-id",
				EndpointID: "endpoint-id",
			}, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "decode cipher text")
		})

		t.Run("With a cipher text that is invalid", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := NewMockStorage(ctrl)

			storedData := EncryptedData{
				Type: EncryptedDataTypeAESCFB,
				Data: "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			}
			storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

			encryptedStorage, err := NewEncryptedStorage(ctx, cfg, storage)
			require.NoError(t, err)

			err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
				ID:         "data-id",
				EndpointID: "endpoint-id",
			}, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid character")
		})

		t.Run("With a valid cipher text", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := NewMockStorage(ctrl)

			storedData := EncryptedData{
				Type: EncryptedDataTypeAESCFB,
				// Hello, World! encrypted with AES CFB and the key
				Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
			}
			storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

			encryptedStorage, err := NewEncryptedStorage(ctx, cfg, storage)
			require.NoError(t, err)

			var decryptedText string
			err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
				ID:         "data-id",
				EndpointID: "endpoint-id",
			}, &decryptedText)
			require.NoError(t, err)
			assert.Equal(t, "Hello, World!", decryptedText)
		})

		t.Run("If the encrypted data link is using the old storage format", func(t *testing.T) {
			encryptedStorage, err := NewEncryptedStorage(ctx, cfg, nil)
			require.NoError(t, err)

			var decryptedText string
			err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
				Type: EncryptedDataTypeAESCFB,
				Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
			}, &decryptedText)
			require.NoError(t, err)
			assert.Equal(t, "Hello, World!", decryptedText)
		})
	})

	t.Run("AES CFB Sha512 decryption", func(t *testing.T) {
		t.Run("With a checksum mismatch", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := NewMockStorage(ctrl)

			storedData := EncryptedData{
				Type: EncryptedDataTypeAESCFBSha512,
				Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
				Hash: "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
			}
			storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

			encryptedStorage, err := NewEncryptedStorage(ctx, cfg, storage)
			require.NoError(t, err)

			err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
				ID:         "data-id",
				EndpointID: "endpoint-id",
			}, nil)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "hash mismatch after decryption")
		})

		t.Run("With a valid cipher text", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			storage := NewMockStorage(ctrl)

			storedData := EncryptedData{
				Type: EncryptedDataTypeAESCFBSha512,
				// Hello, World! encrypted with AES CFB Sha512 and the key
				Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
				// Valid SHA-512 hash of the plaintext "Hello, World!"
				Hash: "625c3af9e72459f50fdff9af15fa7a94b9c589eb1f0a2bca41abd7f6602198bc7ae35bf6c4c296f8039d3af278424500086a783f9b7baa84fad70b41b9e2c6ea",
			}
			storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

			encryptedStorage, err := NewEncryptedStorage(ctx, cfg, storage)
			require.NoError(t, err)

			var decryptedText string
			err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
				ID:         "data-id",
				EndpointID: "endpoint-id",
			}, &decryptedText)
			require.NoError(t, err)
			assert.Equal(t, "Hello, World!", decryptedText)
		})
	})

	t.Run("Decrypt works with alternate keys", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		storage := NewMockStorage(ctrl)

		config := config.Config{
			SecretStorageEncryptionKey: "my-bad-secret-key-is-32-bytes-long-1234567890",
			SecretStorageAlternateKeys: []string{"my-secret-key-is-32-bytes-long-1234567890"},
		}

		storedData := EncryptedData{
			Type: EncryptedDataTypeAESCFBSha512,
			// Hello, World! encrypted with AES CFB Sha512 and the key
			Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
			// Valid SHA-512 hash of the plaintext "Hello, World!"
			Hash: "625c3af9e72459f50fdff9af15fa7a94b9c589eb1f0a2bca41abd7f6602198bc7ae35bf6c4c296f8039d3af278424500086a783f9b7baa84fad70b41b9e2c6ea",
		}
		storage.EXPECT().GetEncryptedData(gomock.Any(), "endpoint-id", "data-id").Return(storedData, nil)

		encryptedStorage, err := NewEncryptedStorage(ctx, config, storage)
		require.NoError(t, err)

		var decryptedText string
		err = encryptedStorage.Decrypt(ctx, EncryptedDataLink{
			ID:         "data-id",
			EndpointID: "endpoint-id",
		}, &decryptedText)
		require.NoError(t, err)
		assert.Equal(t, "Hello, World!", decryptedText)
	})
}

func Test_EncryptedStorage_RotateEncryptionKey(t *testing.T) {
	ctx := context.Background()

	t.Run("with no alternate keys", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		storage := NewMockStorage(ctrl)

		config := config.Config{
			SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
		}

		encryptedStorage, err := NewEncryptedStorage(ctx, config, storage)
		require.NoError(t, err)

		err = encryptedStorage.RotateEncryptionKey(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no alternate keys available for rotation")
	})

	t.Run("with a secret not encrypted with an alternate key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		storage := NewMockStorage(ctrl)

		config := config.Config{
			SecretStorageEncryptionKey: "my-secret-key-is-32-bytes-long-1234567890",
			SecretStorageAlternateKeys: []string{"another-secret-key-is-32-bytes-long-1234567890"},
		}

		storedData := EncryptedData{
			Type: EncryptedDataTypeAESCFBSha512,
			// Hello, World! encrypted with AES CFB Sha512 and the key
			Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
			// Valid SHA-512 hash of the plaintext "Hello, World!"
			Hash:       "625c3af9e72459f50fdff9af15fa7a94b9c589eb1f0a2bca41abd7f6602198bc7ae35bf6c4c296f8039d3af278424500086a783f9b7baa84fad70b41b9e2c6ea",
			ID:         "data-id",
			EndpointID: "endpoint-id",
		}
		storage.EXPECT().ListEncryptedDataForHost(gomock.Any()).Return([]EncryptedData{
			storedData,
		}, nil)

		encryptedStorage, err := NewEncryptedStorage(ctx, config, storage)
		require.NoError(t, err)

		err = encryptedStorage.RotateEncryptionKey(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no alternate key could decrypt the data")
	})

	t.Run("when everything works", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		storage := NewMockStorage(ctrl)

		config := config.Config{
			SecretStorageEncryptionKey: "another-secret-key-is-32-bytes-long-1234567890",
			SecretStorageAlternateKeys: []string{"my-secret-key-is-32-bytes-long-1234567890"},
		}
		encryptedStorage, err := NewEncryptedStorage(ctx, config, storage)

		storedData := EncryptedData{
			Type: EncryptedDataTypeAESCFBSha512,
			// Hello, World! encrypted with AES CFB Sha512 and the key
			Data: "6052a586c3fc190d3865efbd3a7e26404d7f249ef477f2da1db6c71bb6dda0",
			// Valid SHA-512 hash of the plaintext "Hello, World!"
			Hash:       "625c3af9e72459f50fdff9af15fa7a94b9c589eb1f0a2bca41abd7f6602198bc7ae35bf6c4c296f8039d3af278424500086a783f9b7baa84fad70b41b9e2c6ea",
			ID:         "data-id",
			EndpointID: "endpoint-id",
		}
		storage.EXPECT().ListEncryptedDataForHost(gomock.Any()).Return([]EncryptedData{
			storedData,
		}, nil)

		storage.EXPECT().UpsertEncryptedData(gomock.Any(), "endpoint-id", gomock.Any()).DoAndReturn(func(ctx context.Context, endpointID string, data EncryptedData) (EncryptedDataLink, error) {
			// The hash should not have changed
			assert.Equal(t, storedData.Hash, data.Hash)
			// The data should've been re-encrypted with the new key
			assert.NotEqual(t, storedData.Data, data.Data)

			key := sha256.Sum256([]byte(config.SecretStorageEncryptionKey))
			dataBytes, err := hex.DecodeString(data.Data)
			require.NoError(t, err)
			res, err := crypto.Decrypt(key[:], dataBytes)
			require.NoError(t, err)
			assert.Equal(t, "\"Hello, World!\"", string(res))

			assert.Equal(t, EncryptedDataTypeAESCFBSha512, data.Type)
			return EncryptedDataLink{
				ID:         "new-data-id",
				EndpointID: endpointID,
			}, nil
		})
		require.NoError(t, err)
		err = encryptedStorage.RotateEncryptionKey(ctx)
		require.NoError(t, err)
	})
}
