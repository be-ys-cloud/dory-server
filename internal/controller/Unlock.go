package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// Unlock
// @Tags         unlock
// @Summary      Unlock a user.
// @Description  Unlock a user.
// @Param        body           body    structures.UserUnlock  true  "User"
// @Success      200            "OK - User unlocked"
// @Failure      400            "Missing data in payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /unlock [post]
func Unlock(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Missing old password in payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	err = service.UnlockAccount(user.Username, user.Token)

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
