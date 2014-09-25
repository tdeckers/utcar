package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
)

//HttpPost posts item updates to Openhab, based on the provided SIA message.
// Known SIA messages:
// * UA - detector activated
// * UR - detector restored
// * RP - Communication test (e.g. midnight)
// See: http://alarmsbc.com/tech/pdf/sia.pdf
func HttpPost(address string, user string, pwd string, sia SIA) {
	// Openhab: https://asterix.ducbase.com:8443/rest/items/al_{item}/state
	url := strings.Join([]string{"https://", address, "/rest/items/al_", sia.zone, "/state"}, "")
	var body string
	switch sia.command {
	case "UA":
		body = "ON"
	case "UR":
		body = "OFF"
	default:
		log.Printf("Unsupported SIA command for pusher (%s)\n", sia.command)
		return // exit the pusher function
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		log.Panicf("HTTP Request (%v)", err)
	}
	if user != "" && pwd != "" {
		request.SetBasicAuth(user, pwd)
	}
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
