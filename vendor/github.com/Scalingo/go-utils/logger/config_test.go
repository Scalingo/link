package logger

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFormatter(t *testing.T) {
	examples := []struct {
		Name          string
		Env           map[string]string
		FormatterType interface{}
	}{
		{
			Name:          "with no environment",
			FormatterType: &logrus.TextFormatter{},
		}, {
			Name:          "with an invalid environment",
			Env:           map[string]string{"LOGGER_TYPE": "invalid"},
			FormatterType: &logrus.TextFormatter{},
		}, {
			Name:          "with a json type",
			Env:           map[string]string{"LOGGER_TYPE": "json"},
			FormatterType: &logrus.JSONFormatter{},
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			for k, v := range example.Env {
				os.Setenv(k, v)
			}

			assert.IsType(t, example.FormatterType, formatter())
		})
	}
}

func TestLogLevel(t *testing.T) {
	examples := map[string]logrus.Level{
		"panic":   logrus.PanicLevel,
		"fatal":   logrus.FatalLevel,
		"warn":    logrus.WarnLevel,
		"info":    logrus.InfoLevel,
		"debug":   logrus.DebugLevel,
		"":        logrus.InfoLevel,
		"invalid": logrus.InfoLevel,
	}

	for k, v := range examples {
		t.Run(fmt.Sprintf("when LOGGER_LEVEL=%s", k), func(t *testing.T) {
			os.Setenv("LOGGER_LEVEL", k)
			assert.Equal(t, v, logLevel())
		})
	}
}
