package philaehandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Scalingo/go-utils/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PhilaeHandler struct {
	prober *prober.Prober
}

func (handler PhilaeHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	start := time.Now()
	result := handler.prober.Check(ctx)
	json.NewEncoder(response).Encode(result)
	duration := time.Now().Sub(start)

	l := logger.Get(ctx).WithFields(logrus.Fields{
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

func NewHandler(prober *prober.Prober) http.Handler {
	return PhilaeHandler{
		prober: prober,
	}
}

func NewPhilaeRouter(router http.Handler, prober *prober.Prober) *mux.Router {
	globalRouter := mux.NewRouter()
	globalRouter.Handle("/_health", NewHandler(prober))
	globalRouter.Handle("/{any:.+}", router)
	return globalRouter
}
