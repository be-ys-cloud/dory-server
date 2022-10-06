package connectors

import (
	"crypto/tls"
	"github.com/be-ys-cloud/dory-server/internal/structures"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func WSProvider(method string, url string, payload io.Reader, headers map[string]string) (int, []byte, map[string]string, error) {
	res, _ := http.NewRequest(method, url, payload)

	for i, k := range headers {
		res.Header.Set(i, k)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Timeout: time.Second * 20, Transport: tr}

	resp, err := client.Do(res)
	if err != nil {
		return 0, nil, nil, &structures.CustomError{HttpCode: 0, Text: err.Error()}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	headersResp := map[string]string{}

	for name, values := range resp.Header {
		for _, value := range values {
			headersResp[name] = value
		}
	}

	return resp.StatusCode, body, headersResp, nil
}
