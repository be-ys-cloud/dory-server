package controller

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/configuration"
	"github.com/be-ys-cloud/dory-server/internal/service"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/gorilla/mux"
	"net/http"
)

// Ask
// @Tags         demand
// @Summary      Ask server to send email with a link to reset an account password, or unlock it.
// @Description  Ask server to send email with a link to reset an account password, or unlock it.
// @Param        kind           path      string  true  "Kind of request : reinitialize or unlock."
// @Param        body           body    structures.UserAsk  true  "User"
// @Success      200            "OK - Check your mailbox"
// @Failure      400            "Missing username in payload"
// @Failure      500            "An error occurred."
// @Security     BasicAuth
// @Router       /request/{kind} [post]
func Ask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind := vars["kind"]

	// Check kind & feature flipping
	if kind != "unlock" && kind != "reinitialize" {
		w.WriteHeader(404)
		return
	}

	if (kind == "unlock" && configuration.Configuration.Features.DisableUnlock) ||
		(kind == "reinitialize" && configuration.Configuration.Features.DisablePasswordReinitialization) {
		w.WriteHeader(404)
		return
	}

	//Decoding JSON
	var user structures.UserAsk
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" {
		http.Error(w, "Incomplete payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	if err = service.AskMail(user, kind); err != nil {
		handleErrors(err, w)
		return
	}

	w.WriteHeader(200)
}
