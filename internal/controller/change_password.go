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
	var user structures.UserChangePassword
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" || user.OldPassword == "" || user.NewPassword == "" {
		http.Error(w, "Incomplete payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	if err = service.ChangePassword(user); err != nil {
		handleErrors(err, w)
		return
	}

	w.WriteHeader(200)
}
