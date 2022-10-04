package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"net/http"
)

func Config(w http.ResponseWriter, _ *http.Request) {
	d, _ := json.Marshal(configuration.Configuration.Features)
	_, _ = w.Write(d)
}
