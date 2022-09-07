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
	var user structures.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" {
		http.Error(w, "Missing username in payload.", http.StatusBadRequest)
		return
	}

	if user.Token == "" {
		http.Error(w, "Missing token in payload.", http.StatusBadRequest)
		return
	}

	if user.NewPassword == "" {
		http.Error(w, "Missing new password in payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	err = service.ReinitializePassword(user.Username, user.Token, user.NewPassword)

	if err != nil {
		fullError, convertSucceeded := err.(*structures.CustomError)
		if convertSucceeded {
			http.Error(w, fullError.Text, fullError.HttpCode)
			return
		} else {
			http.Error(w, "An unknown (and probably awkward) error occurred.", 500)
			return
		}
	}

	w.WriteHeader(200)

}
