package web

import (
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
)

type versionController struct {
	version string
}

func NewVersionController(version string) versionController {
	return versionController{
		version: version,
	}
}

func (c versionController) Version(w http.ResponseWriter, r *http.Request, params map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"version": c.version,
	})
	if err != nil {
		log.WithError(err).Error("Fail to encode version")
	}
	return nil
}
