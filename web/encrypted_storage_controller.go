package web

import (
	"net/http"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/models"
)

type EncryptedStorageController struct {
	encryptedStorage models.EncryptedStorage
}

func NewEncryptedStorageController(encryptedStorage models.EncryptedStorage) EncryptedStorageController {
	return EncryptedStorageController{
		encryptedStorage: encryptedStorage,
	}
}

func (d EncryptedStorageController) RotateEncryptionKey(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	ctx := r.Context()

	err := d.encryptedStorage.RotateEncryptionKey(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "rotate encryption key")
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
