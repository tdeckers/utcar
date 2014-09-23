package main

import (
	"testing"
	"net"
	"log"
	"sync"
	"io"
)

var (
	serverAddr	string
	once		sync.Once
)

func listen() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:0") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}
	return l, l.Addr().String()
}

func startServer() {
	var l net.Listener
	l, serverAddr = listen()
	log.Println("Test server listening on", serverAddr)
	go acceptConnection(l)
}

func acceptConnection(l net.Listener) {
	conn, e := l.Accept()
	if e != nil {
		log.Fatalf("l.Accept: %v", e)
	}
	handleConnection(conn)	
}

// taking inspiration from http://golang.org/src/pkg/net/rpc/server_test.go
func TestHandleConnection(t *testing.T) {
	once.Do(startServer)
	testHandleConnection(t, serverAddr)
}

func readKey(c net.Conn, t *testing.T) {
        buf := make([]byte, 1024) // receive buffer
        n, err := c.Read(buf)
        if err != nil {
                if err != io.EOF {
                        log.Fatalf("Key read error: ", err)
                }
        }
	if n != 24 {
		t.Errorf("Expected key length is 24, was %v", n)
	}
}

func testHandleConnection(t *testing.T, addr string) {
	client, e := net.Dial("tcp", addr)
	if e != nil {
		t.Fatalf("Dial (%v)", e)
	}
	defer client.Close()
	readKey(client, t)
	// TODO: complete test
}
