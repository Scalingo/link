package philaehandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Scalingo/go-philae/prober"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PhilaeHandler struct {
	prober *prober.Prober
	logger logrus.FieldLogger
}

func (handler PhilaeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	start := time.Now()
	result := handler.prober.Check()
	json.NewEncoder(response).Encode(result)
	duration := time.Now().Sub(start)

	l := handler.logger.WithFields(logrus.Fields{
		"router":   "philae",
		"duration": duration.String(),
		"healthy":  result.Healthy,
	})

	if result.Healthy {
		l.Debug()
	} else {
		l.Info()
	}
}

type HandlerOpts struct {
	Logger logrus.FieldLogger
}

func NewHandler(prober *prober.Prober, opts HandlerOpts) http.Handler {
	h := PhilaeHandler{
		prober: prober,
		logger: opts.Logger,
	}
	if h.logger == nil {
		h.logger = logrus.New()
	}
	return h
}

func NewPhilaeRouter(router http.Handler, prober *prober.Prober, opts HandlerOpts) *mux.Router {
	globalRouter := mux.NewRouter()
	globalRouter.Handle("/_health", NewHandler(prober, opts))
	globalRouter.Handle("/{any:.+}", router)
	return globalRouter
}
