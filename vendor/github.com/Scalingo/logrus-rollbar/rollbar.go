package logrus_rollbar

import (
	"net/http"

	"github.com/rollbar/rollbar-go"
)

type Sender interface {
	RequestError(string, *http.Request, error, map[string]interface{})
	Error(string, error, map[string]interface{})
}

type RollbarSender struct{}

var (
	levelSenders = map[string]func(args ...interface{}){
		rollbar.CRIT:  rollbar.Critical,
		rollbar.ERR:   rollbar.Error,
		rollbar.WARN:  rollbar.Warning,
		rollbar.INFO:  rollbar.Info,
		rollbar.DEBUG: rollbar.Debug,
	}
)

func (s RollbarSender) RequestError(severity string, req *http.Request, err error, fields map[string]interface{}) {
	levelSenders[severity](req, err, fields)
}

func (s RollbarSender) Error(severity string, err error, fields map[string]interface{}) {
	levelSenders[severity](err, fields)
}
