package uaa

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

// UAACheckTokenError is the json structure returned from UAA
// if token is invalid
type UAACheckTokenError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// CheckUAAToken make a post request to uaaCheckTokenURL
// to check whether provided token is valid
func CheckUAAToken(uaaCheckTokenURL, authStr string) (int, error) {
	// Prepare HTTP request
	token := strings.TrimPrefix(authStr, "Bearer ")
	body := []byte("token=" + token)
	req, err := http.NewRequest("POST", uaaCheckTokenURL, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}
	req.SetBasicAuth(os.Getenv("UAA_USERNAME"), os.Getenv("UAA_PASSWORD"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("Accept", "application/json;charset=utf-8")

	// Perform HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	// If status code is not 200, parse & return the error msg
	if resp.StatusCode != http.StatusOK {
		var payload UAACheckTokenError
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			log.Println(err)
		}
		return resp.StatusCode, errors.New(payload.Error)
	}

	return resp.StatusCode, nil
}
