package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"strings"
	"testing"
)

func TestReinitialize(t *testing.T) {

	// Send invalid user to request
	url := baseUrl + "request/reinitialize"

	data := user{
		Username: "unexistent_user",
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, _, _, err := connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 400 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	// Send request for valid user
	url = baseUrl + "request/reinitialize"

	data = user{
		Username: "testuser",
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	// Check mail have been sent and retrieve link
	code, response, _, err := connectors.WSProvider("GET", mailUrl, nil, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	var mails []email
	_ = json.Unmarshal(response, &mails)

	resetToken := strings.Split(strings.Split(mails[0].TextAsHtml, "https://localhost:8001/reinitialize?user=testuser&amp;token=")[1], "\"")[0]

	// Update password to check that token is valid
	// Send request for valid user
	url = baseUrl + "reinitialize"

	data = user{
		Username:       "testuser",
		NewPassword:    "test",
		Authentication: authentication{Token: resetToken},
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	// Same request should end in 400 because token is revoked.
	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 400 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}
}
