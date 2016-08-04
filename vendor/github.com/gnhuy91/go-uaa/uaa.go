package uaa

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
)

// VerifyError is the json structure returned from UAA
// if token is invalid
type VerifyError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// VerifyToken make a post request to uaaURL
// to check whether provided token is valid
func VerifyToken(uaaURL, authStr string) (int, error) {
	// Prepare HTTP request
	token := strings.TrimPrefix(authStr, "Bearer ")
	body := []byte("token=" + token)
	req, err := http.NewRequest("POST", uaaURL, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// Set Basic Auth header using env variables
	// TODO: get username & password from parameters
	req.SetBasicAuth(os.Getenv("UAA_USERNAME"), os.Getenv("UAA_PASSWORD"))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Accept", "application/json;charset=utf-8")

	// Perform HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	// If status code is not 200, parse & return the error msg
	if resp.StatusCode != http.StatusOK {
		var payload VerifyError
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			return resp.StatusCode, err
		}
		return resp.StatusCode, errors.New(payload.Error)
	}

	return resp.StatusCode, nil
}
