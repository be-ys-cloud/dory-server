package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"strings"
	"testing"
)

func TestFailCreateTOTP(t *testing.T) {

	// Trying to create TOTP
	url := baseUrl + "totp/create"

	data := user{
		Username: "testuser",
		Password: "badpassword",
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, _, _, err := connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 401 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}
}
