package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"net/http"
)

// AskReinitialization
// @Tags         reinitialization
// @Summary      Ask server to send email with a link to reset an account password.
// @Description  Ask server to send email with a link to reset an account password.
// @Param        body           body    structures.UserAsk  true  "User (only username is required)"
// @Success      200            "OK - Check your mailbox"
// @Failure      400            "Missing username in payload"
// @Failure      500            "An error occured."
// @Security     BasicAuth
// @Router       /ask_reinitialization [post]
func AskReinitialization(w http.ResponseWriter, r *http.Request) {

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

	//Contact service
	err = service.AskPasswordReinitialization(user.Username)

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
