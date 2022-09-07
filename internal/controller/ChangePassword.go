package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// ChangePassword
// @Tags         change_password
// @Summary      Change a user's password.
// @Description  Change a user's password.
// @Param        body           body    structures.UserChangePassword  true  "User data"
// @Success      200            "OK - Mail changed"
// @Failure      400            "Invalid payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /change_password [post]
func ChangePassword(w http.ResponseWriter, r *http.Request) {

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

	if user.OldPassword == "" {
		http.Error(w, "Missing old password in payload.", http.StatusBadRequest)
		return
	}

	if user.NewPassword == "" {
		http.Error(w, "Missing new password in payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	err = service.ChangePassword(user.Username, user.OldPassword, user.NewPassword)

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
