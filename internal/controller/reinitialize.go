package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// Reinitialize
// @Tags         reinitialization
// @Summary      Reinitialize a user's password.
// @Description  Reinitialize a user's password.
// @Param        body           body    structures.UserReinitialize  true  "User"
// @Success      200            "OK - Password changed"
// @Failure      400            "Missing data in payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /reinitialize [post]
func Reinitialize(w http.ResponseWriter, r *http.Request) {
	//Decoding JSON
	var user structures.UserReinitialize
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" || user.NewPassword == "" {
		http.Error(w, "Incomplete payload.", http.StatusBadRequest)
		return
	}

	if user.Authentication.Token == "" && user.Authentication.TOTP == "" {
		http.Error(w, "Missing authentication method.", http.StatusBadRequest)
		return
	}

	//Contact service
	if err = service.ReinitializePassword(user); err != nil {
		handleErrors(err, w)
		return
	}

	w.WriteHeader(200)

}
