package mivo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	authKey              string = ""
	currentSign          string = ""
	currentSignLastFetch int64  = 0
)

type wmsAuth struct {
	Sign string
}

func getAuth() (string, error) {
	if authKey != "" {
		return authKey, nil
	}

	file, err := os.Open("auth.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GetSign() (string, error) {
	now := time.Now().Unix()
	if currentSign != "" && currentSignLastFetch+600 > now {
		return currentSign, nil
	}

	_authKey, err := getAuth()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("GET", "https://api.mivo.com/v4/web/channels/wms-auth", nil)
	if err != nil {
		return "", err
	}
	setHeaders(&req.Header)
	req.Header.Set("Authorization", _authKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	dec := json.NewDecoder(resp.Body)
	var auth wmsAuth
	dec.Decode(&auth)
	return auth.Sign, nil
}
