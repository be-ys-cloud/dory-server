package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"strings"
	"testing"
)

func TestDisabledFeature(t *testing.T) {

	url := baseUrl + "request/unlock"

	data := user{
		Username: "testuser",
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, _, _, err := connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 404 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	url = baseUrl + "unlock"

	data = user{
		Username:       "testuser",
		Authentication: authentication{Token: "test"},
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 404 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}
}
