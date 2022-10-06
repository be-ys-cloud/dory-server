package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"strings"
	"testing"
)

func TestChangePassword(t *testing.T) {

	// Trying to create TOTP
	url := baseUrl + "change_password"

	data := user{
		Username:    "testuser",
		OldPassword: "badpassword",
		NewPassword: "newpassword",
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, _, _, err := connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 401 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	data = user{
		Username:    "testuser",
		OldPassword: "test",
		NewPassword: "newpassword",
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	data = user{
		Username:    "testuser",
		OldPassword: "newpassword",
		NewPassword: "test",
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

}
