package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpPost(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Panicf("Request body read (%v)", e)
		}
		if r.Method != "POST" {
			t.Errorf("Reqeust method isn't POST, but %s", r.Method)
		}
		if string(body) != "ON" {
			t.Errorf("Body - expected ON, was %s", body)
		}
	})

	server := httptest.NewTLSServer(handler)
	defer server.Close()
	address := server.Listener.Addr().String()
	log.Printf("HTTP test server (%v)\n", address)
	sia := SIA{time.Now(), "", "", "", "", "UA", "272"}
	HttpPost(address, "tom", "pwd", sia)
}
