package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//HttpPost posts item updates to Openhab, based on the provided SIA message.
// Known SIA messages:
// * UA - detector activated
// * UR - detector restored
// * RP - Communication test (e.g. midnight)
// See: http://alarmsbc.com/tech/pdf/sia.pdf
func HttpPost(address string, user string, pwd string, sia SIA) error {
	// Parse base URL
	_, err := url.Parse(address)
	if err != nil {
		log.Panicf("Address parse failed (%v)", err)
	}
	// Openhab: https://asterix.ducbase.com:8443/rest/items/al_{item}/state
	itemUrl := strings.Join([]string{address, "/rest/items/al_", sia.zone, "/state"}, "")
	var body string
	switch sia.command {
	case "UA":
		body = "ON"
	case "UR":
		body = "OFF"
	default:
		return fmt.Errorf("Unsupported SIA command for pusher (%s)\n", sia.command)
	}
	request, err := http.NewRequest("PUT", itemUrl, strings.NewReader(body))
	if err != nil {
		log.Panicf("HTTP Request (%v)", err)
	}
	request.Header.Add("Content-Type", "text/plain")
	if user != "" && pwd != "" {
		request.SetBasicAuth(user, pwd)
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("HTTP Response (%v)", err)
	}
	defer response.Body.Close()
	log.Printf("PUT %s to %s (%s)", body, itemUrl, response.Status)
	return nil
}
