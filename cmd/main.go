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
	//Get configuration and propagate it to other packages that need it

	m := mux.NewRouter()
	m.Use(controller.SetHeadersMiddleware)

	api := m.PathPrefix(configuration.Configuration.Server.BasePath).Subrouter()

	api.HandleFunc("/ask_reinitialization", controller.AskReinitialization).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/reinitialize", controller.Reinitialize).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/ask_unlock", controller.AskUnlock).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/unlock", controller.Unlock).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/change_password", controller.ChangePassword).Methods(http.MethodPost, http.MethodOptions)

	logrus.Info("Web server started. Now listening on *:" + strconv.Itoa(configuration.Configuration.Server.Port))
	logrus.Fatalln(http.ListenAndServe(":"+strconv.Itoa(configuration.Configuration.Server.Port), m))
}
