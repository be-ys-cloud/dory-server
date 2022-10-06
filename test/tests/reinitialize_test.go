package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"strings"
	"testing"
)

func TestReinitialize(t *testing.T) {

	// Send invalid user to request
	url := baseUrl + "ask/reinitialize"

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
	url = baseUrl + "ask/reinitialize"

	data = user{
		Username: "unexistent_user",
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
