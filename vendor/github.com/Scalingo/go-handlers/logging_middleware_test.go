package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingMiddleware_Apply(t *testing.T) {
	examples := []struct {
		Name           string
		Filters        map[string]logrus.Level
		Handler        func(t *testing.T) HandlerFunc
		Path           string
		Method         string
		Host           string
		Headers        map[string]string
		Context        func(context.Context) context.Context
		Env            map[string]string
		ExpectedLevel  logrus.Level
		ExpectedFields []string
	}{
		{
			Name:           "HTTP GET / on example.dev without any additional info",
			Path:           "/",
			Method:         "GET",
			Host:           "example.dev",
			ExpectedFields: []string{"path", "host", "method"},
		}, {
			Name:           "with user agent",
			Path:           "/",
			Method:         "GET",
			Host:           "example.dev",
			Headers:        map[string]string{"User-Agent": "MyUserAgent1.0"},
			ExpectedFields: []string{"path", "host", "method", "user_agent"},
		}, {
			Name:           "with referer",
			Path:           "/",
			Method:         "GET",
			Host:           "example.dev",
			Headers:        map[string]string{"Referer": "http://google.com"},
			ExpectedFields: []string{"path", "host", "method", "referer"},
		}, {
			Name:           "with request_id in context",
			Path:           "/",
			Method:         "GET",
			Host:           "example.dev",
			ExpectedFields: []string{"path", "host", "method", "request_id"},
			Context: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, "request_id", "0")
			},
		}, {
			Name:           "with request_id in context",
			Path:           "/",
			Method:         "GET",
			Host:           "example.dev",
			ExpectedFields: []string{"path", "host", "method"},
			Handler: func(t *testing.T) HandlerFunc {
				return HandlerFunc(func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
					logger, ok := r.Context().Value("logger").(logrus.FieldLogger)
					assert.True(t, ok)
					assert.NotNil(t, logger)
					return nil
				})
			},
		}, {
			Name:           "with filtered path it should not do anything if feature not enabled",
			Path:           "/res/1/health",
			Filters:        map[string]logrus.Level{"/res/.+/health": logrus.DebugLevel},
			Method:         "GET",
			Host:           "example.dev",
			ExpectedLevel:  logrus.InfoLevel,
			ExpectedFields: []string{"path", "host", "method"},
		}, {
			Name:           "with filtered path it should change the level if the feature is enabled",
			Path:           "/res/1/health",
			Filters:        map[string]logrus.Level{"/res/.+/health": logrus.DebugLevel},
			Method:         "GET",
			Host:           "example.dev",
			Env:            map[string]string{"HANDLERS_LOG_FILTERS": "true"},
			ExpectedLevel:  logrus.DebugLevel,
			ExpectedFields: []string{"path", "host", "method"},
		}, {
			Name:           "should not filter unfiltered path",
			Path:           "/res/health",
			Filters:        map[string]logrus.Level{"/res/.+/health": logrus.DebugLevel},
			Method:         "GET",
			Host:           "example.dev",
			ExpectedLevel:  logrus.InfoLevel,
			ExpectedFields: []string{"path", "host", "method"},
		},
	}

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request, params map[string]string) error {
		return nil
	})

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			var err error

			envbackup := map[string]string{}
			if example.Env != nil {
				for k, v := range example.Env {
					envbackup[k] = os.Getenv(k)
					os.Setenv(k, v)
				}
				defer func() {
					for k, v := range envbackup {
						os.Setenv(k, v)
					}
				}()
			}

			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.DebugLevel)
			defer hook.Reset()

			middleware := NewLoggingMiddleware(logger)
			if example.Filters != nil {
				middleware, err = NewLoggingMiddlewareWithFilters(logger, example.Filters)
				require.NoError(t, err)
			}

			reqHandler := handler
			if example.Handler != nil {
				reqHandler = example.Handler(t)
			}
			reqHandler = middleware.Apply(reqHandler)

			w := httptest.NewRecorder()
			r, err := http.NewRequest(example.Method, example.Path, nil)
			require.NoError(t, err)

			if example.Context != nil {
				r = r.WithContext(example.Context(r.Context()))
			}

			r.Host = example.Host
			if example.Headers != nil {
				for k, v := range example.Headers {
					r.Header.Add(k, v)
				}
			}

			err = reqHandler(w, r, map[string]string{})
			require.NoError(t, err)

			require.Equal(t, 2, len(hook.Entries))

			expectedLevel := logrus.InfoLevel
			if example.ExpectedLevel != 0 {
				expectedLevel = example.ExpectedLevel
			}

			assert.Equal(t, expectedLevel, hook.Entries[0].Level)
			assert.Equal(t, expectedLevel, hook.Entries[1].Level)

			log1Keys := []string{}
			for k, _ := range hook.Entries[0].Data {
				log1Keys = append(log1Keys, k)
			}
			log2Keys := []string{}
			for k, _ := range hook.Entries[1].Data {
				log2Keys = append(log2Keys, k)
			}

			assert.Subset(t, log1Keys, example.ExpectedFields)
			assert.Subset(t, log1Keys, []string{"protocol"})

			assert.Subset(t, log2Keys, example.ExpectedFields)
			assert.Subset(t, log2Keys, []string{"protocol", "status", "duration", "bytes"})
		})
	}
}
