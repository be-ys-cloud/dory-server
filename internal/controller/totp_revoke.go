package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// RevokeTOTP
// @Tags         totp
// @Summary      Revoke all TOTP tokens for this user.
// @Description  Revoke all TOTP tokens for this user.
// @Param        body           body    structures.UserCreateTOTP  true  "User data"
// @Success      200            "TOTP deleted"
// @Failure      400            "Invalid payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /totp/create [post]
func RevokeTOTP(w http.ResponseWriter, r *http.Request) {

	//Decoding JSON
	var user structures.UserCreateTOTP
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" || user.Password == "" {
		http.Error(w, "Incomplete payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	if err = service.RevokeTOTP(user); err != nil {
		handleErrors(err, w)
		return
	}
}
