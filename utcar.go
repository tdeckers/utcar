package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type SIA struct {
	time     time.Time
	sequence string
	receiver string
	line     string
	account  string
	command  string
	zone     string
}

type Heartbeat struct {
	time time.Time
}

var ftaddr string
var ftuser string
var ftpwd string
var fport int

// init function.  Used to read input parameters to the program.
func init() {
	flag.StringVar(&ftaddr, "taddr", "", "Target addr (host:port)")
	flag.StringVar(&ftuser, "tuser", "", "Target username")
	flag.StringVar(&ftpwd, "tpwd", "", "Target password")
	flag.IntVar(&fport, "port", 12300, "Listen port number (default: 12300)")
	flag.Parse()
}

// handleConnection handles connections from the alarm system.
// In short, it accepts a connection and sends a new, encrypted key.  Then it
// receives an encrypted message from the alarm system, after which it completes
// with an ACK message.
func handleConnection(c net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Message processing panic (%v)\n", r)
		}
	}()
	key := GenerateKey()
	scrambled_key := Scramble(key)
	// Send key to alarm system
	n, err := c.Write(scrambled_key)
	if err != nil {
		log.Panic(err)
	}

	buf := make([]byte, 1024) // receive buffer
	n, err = c.Read(buf)
	if err != nil {
		if err != io.EOF {
			log.Panic("Read error: ", err)
		}
	}
	encryptedData := buf[:n]

	data := Decrypt3DESECB(encryptedData, key)
	log.Println("Message: ", string(data[:]))

	ack := []byte("ACK\r")
	ack = append(ack, []byte{0, 0, 0, 0}...)
	encryptedAck := Encrypt3DESECB(ack, key)
	n, err = c.Write(encryptedAck)
	if err != nil {
		log.Panic(err)
	}

	if IsHeartbeat(data) {
		return // don't know what to do with this yet.
	}
	parsed := ParseSIA(data)
	if parsed == nil {
		log.Panicf("Not a recognized message: %s", string(data[:]))
	}
	sia := SIA{time.Now(), parsed[0], parsed[1], parsed[2], parsed[3], parsed[4], parsed[5]}
	log.Println(sia)

}

func main() {
	// Listen on TCP port 12300 on all interfaces
	l, err := net.Listen("tcp", ":"+strconv.Itoa(fport))
	if err != nil {
		log.Fatal(err) // exit.. something serious must be wrong.
	}
	log.Printf("Listing on port %d...", fport)
	defer l.Close()
	for { // eternally...
		// Wait for a connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new routine
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			defer c.Close()

			// Generate a random key
			key := make([]byte, 24)
			rand.Read(key)
			if err != nil {
				log.Fatal(err)
			}
			scrambled_key := Scramble(key)
			log.Printf("Key: %s", hex.EncodeToString(key))
			//log.Printf("Scrambled key: %s", hex.EncodeToString(scrambled_key))
			// Send key to alarm system
			n, err := c.Write(scrambled_key)
			// TODO: compare n with size of key
			if err != nil {
				log.Fatal(err)
				return
			}
			//log.Printf("Sent %d bytes to alarm (key)", n)
			buf := make([]byte, 1024)
			n, err = c.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Fatal("Read error: ", err)
				}
			}
			//log.Printf("Read %d bytes", n)
			encryptedData := buf[:n]
			//log.Printf("Data: %s", hex.EncodeToString(encryptedData))
			data := Decrypt3DESECB(encryptedData, key)
			//fmt.Println("Message(byte): ", hex.EncodeToString(data))
			log.Println("Message: ", string(data[:]))
			ack := []byte("ACK\r")
			ack = append(ack, []byte{0, 0, 0, 0}...)
			encryptedAck := Encrypt3DESECB(ack, key)
			//log.Printf("Encrypted ACK: %s", hex.EncodeToString(encryptedAck))
			n, err = c.Write(encryptedAck)
			if err != nil {
				log.Fatal(err)
			}

		}(conn)
	}
}
