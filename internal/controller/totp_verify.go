package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// VerifyTOTP
// @Tags         totp
// @Summary      Verify a TOTP token for this user.
// @Description  Verify a TOTP for this user.
// @Param        body           body    structures.UserVerifyTOTP  true  "User data"
// @Success      200            "TOTP valid"
// @Success      401            "TOTP invalid"
// @Failure      400            "Invalid payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /totp/create [post]
func VerifyTOTP(w http.ResponseWriter, r *http.Request) {

	//Decoding JSON
	var user structures.UserVerifyTOTP
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" || user.TOTP == "" {
		http.Error(w, "Incomplete payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	result, err := service.CheckTOTP(user)

	if err != nil {
		handleErrors(err, w)
		return
	}

	if result {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(401)
	}
}
