package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ebuildy/mattermost-plugin-minotor/server/config"
	// "github.com/ebuildy/mattermost-plugin-minotor/server/registry"

	"github.com/mattermost/mattermost/server/public/pluginapi"
)

type Handler struct {
	*ErrorHandler
	pluginAPI *pluginapi.Client
	config    config.Service
}

func NewHandler(pluginAPI *pluginapi.Client, config config.Service) *Handler {
	handler := &Handler{
		ErrorHandler: &ErrorHandler{},
		pluginAPI:    pluginAPI,
		config:       config,
	}

	return handler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// registry := registry.NewRegistry()
}

// handleResponseWithCode logs the internal error and sends the public facing error
// message as JSON in a response with the provided code.
func handleResponseWithCode(w http.ResponseWriter, code int, publicMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	responseMsg, _ := json.Marshal(struct {
		Error string `json:"error"` // A public facing message providing details about the error.
	}{
		Error: publicMsg,
	})
	_, _ = w.Write(responseMsg)
}

// HandleErrorWithCode logs the internal error and sends the public facing error
// message as JSON in a response with the provided code.
func HandleErrorWithCode(logger logrus.FieldLogger, w http.ResponseWriter, code int, publicErrorMsg string, internalErr error) {
	if internalErr != nil {
		logger = logger.WithError(internalErr)
	}

	if code >= http.StatusInternalServerError {
		logger.Error(publicErrorMsg)
	} else {
		logger.Warn(publicErrorMsg)
	}

	handleResponseWithCode(w, code, publicErrorMsg)
}