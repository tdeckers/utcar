package main

import (
	"bytes"
	"io"
	"log"
	"net"
	"sync"
	"testing"
)

var (
	serverAddr string
	once       sync.Once
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
	handleConnection(conn, nil)
}

// taking inspiration from http://golang.org/src/pkg/net/rpc/server_test.go
func TestHandleConnection(t *testing.T) {
	once.Do(startServer)
	testHandleConnection(t, serverAddr)
}

func readKey(c net.Conn, t *testing.T) []byte {
	buf := make([]byte, 1024) // receive buffer
	n, err := c.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Fatalf("Key read error: %v", err)
		}
	}
	if n != 24 {
		t.Errorf("Expected key length is 24, was %v", n)
	}
	return Scramble(buf[:n])
}

func testHandleConnection(t *testing.T, addr string) {
	client, e := net.Dial("tcp", addr)
	if e != nil {
		t.Fatalf("Dial (%v)", e)
	}
	defer client.Close()
	key := readKey(client, t)
	data := []byte("01010053\"SIA-DCS\"0007R0075L0001[#001465|NRP000*'DECKERS'NM]7C9677F21948CC12|#001465")
	data = append(data, []byte{0, 0, 0, 0, 0}...)
	encrypted := Encrypt3DESECB(data, key)
	_, e = client.Write(encrypted)
	if e != nil {
		t.Fatalf("Write encrypted (%v)", e)
	}
	buf := make([]byte, 16) // only need 8
	n, e := client.Read(buf)
	if e != nil {
		t.Fatalf("Failed to read ACK (%v)", e)
	}
	if n != 8 {
		t.Fatalf("Expected 8 bytes, read %d", n)
	}
	ack := Decrypt3DESECB(buf[:8], key)
	valid := []byte("ACK\r")
	valid = append(valid, []byte{0, 0, 0, 0}...)
	if !bytes.Equal(valid, ack) {
		t.Fatalf("ACK messages didn't match, was %v", ack)
	}
}
