package main

import (
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/controller"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func main() {

	m := mux.NewRouter()
	m.Use(controller.SetHeadersMiddleware)

	api := m.PathPrefix(configuration.Configuration.Server.BasePath).Subrouter()

	api.HandleFunc("/request/{kind}", controller.Ask).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/config", controller.Config).Methods(http.MethodGet, http.MethodOptions)

	if !configuration.Configuration.Features.DisablePasswordReinitialization {
		api.HandleFunc("/reinitialize", controller.Reinitialize).Methods(http.MethodPost, http.MethodOptions)
	}

	if !configuration.Configuration.Features.DisableUnlock {
		api.HandleFunc("/unlock", controller.Unlock).Methods(http.MethodPost, http.MethodOptions)
	}

	if !configuration.Configuration.Features.DisablePasswordUpdate {
		api.HandleFunc("/change_password", controller.ChangePassword).Methods(http.MethodPost, http.MethodOptions)
	}

	if !configuration.Configuration.Features.DisableTOTP {
		api.HandleFunc("/totp/create", controller.CreateTOTP).Methods(http.MethodPost, http.MethodOptions)
		api.HandleFunc("/totp/verify", controller.VerifyTOTP).Methods(http.MethodPost, http.MethodOptions)
		api.HandleFunc("/totp/revoke", controller.RevokeTOTP).Methods(http.MethodPost, http.MethodOptions)
	}

	logrus.Info("Web server started. Now listening on *:" + strconv.Itoa(configuration.Configuration.Server.Port))
	logrus.Fatalln(http.ListenAndServe(":"+strconv.Itoa(configuration.Configuration.Server.Port), m))
}
