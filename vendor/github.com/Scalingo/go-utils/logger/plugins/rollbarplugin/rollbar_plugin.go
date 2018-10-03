package rollbarplugin

import (
	"os"

	"github.com/Scalingo/go-utils/logger"
	logrus_rollbar "github.com/Scalingo/logrus-rollbar"
	"github.com/rollbar/rollbar-go"
	"github.com/sirupsen/logrus"
)

type RollbarPlugin struct{}

// Register the plugin to the logger library
func Register() {
	logger.Plugins().RegisterPlugin(RollbarPlugin{})
}

func (p RollbarPlugin) Name() string {
	return "rollbar"
}

// Generate the hook
func (p RollbarPlugin) Hook() (bool, logrus.Hook) {
	token := os.Getenv("ROLLBAR_TOKEN")

	if token == "" {
		return false, nil
	}

	rollbar.SetToken(token)
	rollbar.SetEnvironment(os.Getenv("GO_ENV"))

	return true, logrus_rollbar.New()
}
