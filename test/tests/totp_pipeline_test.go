package tests

import (
	"encoding/json"
	"github.com/be-ys-cloud/dory-server/test/connectors"
	"github.com/pquerna/otp/totp"
	"strings"
	"testing"
	"time"
)

// TestTOTP generates a TOTP token, verifies it and revoke it.
func TestTOTP(t *testing.T) {

	// Trying to create TOTP
	url := baseUrl + "totp/create"

	data := user{
		Username: "testuser",
		Password: "test",
	}

	marshaled, _ := json.Marshal(data)
	reader := strings.NewReader(string(marshaled))

	code, response, _, err := connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	var TOTP totpStruct
	err = json.Unmarshal(response, &TOTP)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	TOTP.TOTP = strings.Split(TOTP.TOTP, "secret=")[1]
	TOTP.TOTP = strings.Split(TOTP.TOTP, "&")[0]
	t.Log(TOTP.TOTP)

	// Verifying TOTP
	totpcode, err := totp.GenerateCode(TOTP.TOTP, time.Now())
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	url = baseUrl + "totp/verify"

	data = user{
		Username: "testuser",
		TOTP:     totpcode,
	}

	marshaled, _ = json.Marshal(data)
	reader = strings.NewReader(string(marshaled))

	code, _, _, err = connectors.WSProvider("POST", url, reader, nil)
	if err != nil || code != 200 {
		t.Log(err)
		t.Log(code)
		t.FailNow()
	}

	// Revoking TOTP
	url = baseUrl + "totp/revoke"

	data = user{
		Username: "testuser",
		Password: "test",
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
