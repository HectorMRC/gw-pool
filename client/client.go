package client

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var url = "http://localhost:8080"
var ctype = "application/json"

// SetURL sets the url whero to send all requests
func SetURL(newURL string) {
	url = newURL
}

// PostRequest makes a standard request for gw-pool service
func PostRequest(latitude, longitude, clientID int) (err error) {
	var body []byte
	if body, err = json.Marshal(map[string]int{
		"latitude":  latitude,
		"longitude": longitude,
		"client_id": clientID,
	}); err != nil {
		return
	}

	buffer := bytes.NewBuffer(body)
	var resp *http.Response
	if resp, err = http.Post(url, ctype, buffer); err != nil {
		return
	}

	resp.Body.Close()
	return
}
