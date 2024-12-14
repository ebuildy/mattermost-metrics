package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/ebuildy/mattermost-plugin-minotor/server/exporter"
)

func TestServeHTTP(t *testing.T) {
	assert := assert.New(t)
	plugin := Plugin{
		router:   mux.NewRouter(),
		registry: exporter.NewRegistry(),
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exporter", nil)

	plugin.router.Handle("/metrics", plugin.registry.HandleHTTP())

	plugin.ServeHTTP(nil, w, r)

	result := w.Result()
	assert.NotNil(result)
	defer result.Body.Close()
	bodyBytes, err := io.ReadAll(result.Body)
	assert.Nil(err)
	bodyString := string(bodyBytes)

	assert.Contains(bodyString, "go_gc_duration_seconds")
}
