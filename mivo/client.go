package mivo

import (
	"net/http"
)

func setHeaders(head *http.Header) {
	head.Set("Referer", "https://www.mivo.com/")
	head.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:50.0) Gecko/20100101 Firefox/50.0")
	head.Set("Accept", "application/json")
}

func httpGet(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(&req.Header)
	client := &http.Client{}
	return client.Do(req)
}
