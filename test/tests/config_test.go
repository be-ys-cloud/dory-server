package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"testing"
)

func TestConfig(t *testing.T) {

	// Trying to create TOTP
	url := baseUrl + "config"

	code, response, _, err := connectors.WSProvider("GET", url, nil, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	var conf structures.Features
	_ = json.Unmarshal(response, &conf)

	if conf.DisablePasswordReinitialization || conf.DisableTOTP || !conf.DisableUnlock || conf.DisablePasswordUpdate {
		t.Log("Bad conf")
		t.FailNow()
	}

}
