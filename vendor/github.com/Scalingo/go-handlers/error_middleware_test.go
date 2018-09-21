package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	errorutils "github.com/Scalingo/go-utils/errors"
	"github.com/Scalingo/go-utils/logger"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestErrorMiddlware(t *testing.T) {
	handler := ErrorMiddleware(HandlerFunc(func(w http.ResponseWriter, r *http.Request, vars map[string]string) error {
		log := logger.Get(r.Context()).WithField("field", "value")
		return errorutils.Wrapf(logger.ToCtx(context.Background(), log), errors.New("error"), "wrapping")
	}))

	log, hook := test.NewNullLogger()
	log.SetLevel(logrus.DebugLevel)
	defer hook.Reset()

	ctx := logger.ToCtx(context.Background(), log)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	err := handler(w, r, map[string]string{})
	require.Error(t, err)

	require.Equal(t, 1, len(hook.Entries))
	require.Equal(t, "value", hook.Entries[0].Data["field"])
}
