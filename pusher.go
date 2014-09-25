package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
)

func HttpPost(address string, user string, pwd string, sia SIA) {
	// Openhab: https://asterix.ducbase.com:8443/rest/items/{item}/state
	url := strings.Join([]string{"https://", address, "/rest/items/", sia.zone, "/state"}, "")
	var body string
	switch sia.command {
	case "UA":
		body = "ON"
		break
	case "UR":
		body = "OFF"
		break
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		log.Panicf("HTTP Request (%v)", err)
	}
	request.SetBasicAuth(user, pwd)
	log.Printf("About to POST to %s\n", url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		log.Panicf("HTTP Response (%v)", err)
	}
	defer response.Body.Close()
	log.Println(response)
}
