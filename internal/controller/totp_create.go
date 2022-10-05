package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// CreateTOTP
// @Tags         totp
// @Summary      Create a TOTP token for this user.
// @Description  Create a TOTP for this user.
// @Param        body           body    structures.UserCreateTOTP  true  "User data"
// @Success      201            {object} structures.TOTPToken "TOTP created"
// @Failure      400            "Invalid payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /totp/create [post]
func CreateTOTP(w http.ResponseWriter, r *http.Request) {

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
	data, err := service.CreateTOTP(user)

	if err != nil {
		handleErrors(err, w)
		return
	}

	w.WriteHeader(200)

	d, _ := json.Marshal(data)
	_, _ = w.Write(d)

}
