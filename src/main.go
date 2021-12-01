package main

import (
	"ad"
	"configuration"
	"controller"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"mailer"
	"net/http"
	"service"
	"strconv"
)

func main() {
	//Get configuration and propagate it to other packages that need it
	conf := configuration.GetConfiguration()
	ad.Conf = conf
	mailer.Conf = conf
	service.Conf = conf


	m := mux.NewRouter()

	api := m.PathPrefix(conf.Server.BasePath).Subrouter()

	api.HandleFunc("/ask_reinitialization", controller.AskReinitialization).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/reinitialize", controller.Reinitialize).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/ask_unlock", controller.AskUnlock).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/unlock", controller.Unlock).Methods(http.MethodPost, http.MethodOptions)

	api.HandleFunc("/change_password", controller.ChangePassword).Methods(http.MethodPost, http.MethodOptions)

	logrus.Info("Web server started. Now listening on *:"+strconv.Itoa(conf.Server.Port))
	logrus.Fatalln(http.ListenAndServe(":"+strconv.Itoa(conf.Server.Port), m))
}
