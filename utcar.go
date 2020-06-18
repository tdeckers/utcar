package main

import (
	"bytes"
	"expvar"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

var (
	pchan chan SIA

	requests = expvar.NewInt("requests")
)

var rootCmd = &cobra.Command{
	Use:   "utcar",
	Short: "Utcar provides integration for ATS2000IP alarm system",
	Long: `Utcar provides integration for ATS2000IP alarm system
			and optionally posts to an Openhab home automation
			system.
			Complete documentation is available at 
			https://github.com/tdeckers/utcar`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func Execute() {
	rootCmd.PersistentFlags().String("addr", "", "Target addr (e.g. http://openhab.local:8080)")
	rootCmd.PersistentFlags().String("user", "", "Target username")
	rootCmd.PersistentFlags().String("pwd", "", "Target password")
	rootCmd.PersistentFlags().String("mqttaddr", "", "Mqtt addr (e.g. http://mqtt:1883)")
	rootCmd.PersistentFlags().String("mqttuser", "", "Mqtt username")
	rootCmd.PersistentFlags().String("mqttpwd", "", "Mqtt password")
	rootCmd.PersistentFlags().Int("port", 12300, "Listen port number")
	rootCmd.PersistentFlags().Int("debug", 0, "Debug server port number (default: no debug server)")
	viper.BindPFlag("addr", rootCmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("pwd", rootCmd.PersistentFlags().Lookup("pwd"))
	viper.BindPFlag("mqttaddr", rootCmd.PersistentFlags().Lookup("mqttaddr"))
	viper.BindPFlag("mqttuser", rootCmd.PersistentFlags().Lookup("mqttuser"))
	viper.BindPFlag("mqttpwd", rootCmd.PersistentFlags().Lookup("mqttpwd"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix("utcar") // uppercased automatically
	viper.AutomaticEnv()
	viper.SetDefault("port", 12300)
}

// handleConnection handles connections from the alarm system.
// In short, it accepts a connection and sends a new, encrypted key.  Then it
// receives an encrypted message from the alarm system, after which it completes
// with an ACK message.
func handleConnection(c net.Conn, q chan SIA) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Message processing panic (%v)\n", r)
			debug.PrintStack()
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
	// Remove leading/trailing new line, line feeds, NUL chars
	data = bytes.Trim(data, "\n\r\x00")
	log.Println("Message: ", string(data[:]))

	ack := []byte("ACK\r")
	ack = append(ack, []byte{0, 0, 0, 0}...)
	encryptedAck := Encrypt3DESECB(ack, key)
	n, err = c.Write(encryptedAck)
	if err != nil {
		log.Panic(err)
	}

	if IsHeartbeat(data) {
		log.Println("Heartbeat.")
		return // don't know what to do with this yet.
	}
	parsed, err := ParseSIA(data)
	if err != nil {
		log.Panicf("Not a recognized message: %s", string(data[:]))
	}
	sia := SIA{time.Now(), parsed[0], parsed[1], parsed[2], parsed[3], parsed[4], parsed[5]}
	log.Println(sia)

	requests.Add(1) // accessible through expvar

	if q == nil {
		return
	} else {
		q <- sia
	}
}

func receiveSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		os.Exit(0)
	}()
}

func run() {
	// setup response to CTRL-C
	receiveSignal()
	// Listen on TCP port 12300 on all interfaces
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("port")))
	if err != nil {
		log.Fatal(err) // exit.. something serious must be wrong.
	}
	log.Printf("Listing on port %d...", viper.GetInt("port"))
	defer l.Close()

	// setup debug server
	if viper.GetInt("debug") != 0 {
		go func() {
			err = http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("debug")), nil)
		}()
		if err != nil {
			log.Printf("Failed to start debug server (%v)\n", err)
		} else {
			log.Printf("Debug server running on port %d\n", viper.GetInt("debug"))
		}
	}

	// setup pusher channel (if addr is provided)
	if viper.GetString("addr") != "" {
		log.Printf("Pushing to %s\n", viper.GetString("addr"))
		pchan = make(chan SIA)
		go func() {
			for {
				sia := <-pchan
				// TODO: handle panics from this function (if any?)
				err := HttpPost(viper.GetString("addr"), viper.GetString("user"), viper.GetString("pwd"), sia)
				if err != nil {
					log.Printf("Push error: %v", err)
				}
			}
		}()
	}
	if viper.GetString("mqttaddr") != "" {
		log.Printf("Pushing to mqtt %s\n", viper.GetString("mqttaddr"))
		pchan = make(chan SIA)
		go func() {
			client := ConnectMqtt(viper.GetString("mqttaddr"), viper.GetString("mqttuser"), viper.GetString("mqttpwd"))
			for {
				sia := <-pchan
				err := PublishMqtt(client, sia)
				if err != nil {
					log.Printf("Push error: %v", err)
				}
			}
		}()
	}

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

			handleConnection(c, pchan)
		}(conn)
	}
}

func main() {
	Execute()
}
