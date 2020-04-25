package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpPost(t *testing.T) {
	sia := SIA{time.Now(), "", "", "", "", "UA", "272"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, e := ioutil.ReadAll(r.Body)
		if e != nil {
			log.Panicf("Request body read (%v)", e)
		}
		if r.Method != "POST" {
			t.Errorf("Reqeust method isn't POST, but %s", r.Method)
		}
		if r.URL.Path != "/rest/items/al_"+sia.zone+"/state" {
			t.Errorf("Request Path isn't right (%s)", r.URL.Path)
		}
		if string(body) != "ON" {
			t.Errorf("Body - expected ON, was %s", body)
		}
	})

	server := httptest.NewTLSServer(handler)
	defer server.Close()
	address := fmt.Sprintf("http://%s", server.Listener.Addr().String())
	log.Printf("HTTP test server (%v)\n", address)
	HttpPost(address, "tom", "pwd", sia)
}
