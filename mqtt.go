package main

import (
	"crypto/tls"
	//"flag"
	"fmt"
	"log"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)


func ConnectMqtt(mqttaddress string, mqttuser string, mqttpwd string) MQTT.Client {
	//MQTT.DEBUG = log.New(os.Stdout, "", 0)
	//MQTT.ERROR = log.New(os.Stdout, "", 0)
	//hostname, _ := os.Hostname()
	//topic := flag.String("topic", "alarm/", "Topic to publish the messages on")
	//qos := flag.Int("qos", 0, "The QoS to send the messages at")
	//retained := flag.Bool("retained", false, "Are the messages sent with the retained flag")
	//clientid := flag.String("clientid", "utcar", "A clientid for the connection")
	//flag.Parse()
	
	connOpts := MQTT.NewClientOptions().AddBroker(mqttaddress).SetClientID("utcar").SetCleanSession(true)
	connOpts.SetUsername(mqttuser)
	connOpts.SetPassword(mqttpwd)
	
	tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
	connOpts.SetTLSConfig(tlsConfig)

	client := MQTT.NewClient(connOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return nil
	}
	log.Printf("Connected to %s\n", mqttaddress)
	return client
}

func PublishMqtt(client MQTT.Client, sia SIA) error {
	var body string
	switch sia.command {
	case "UA":
		body = "ON"
	case "UR":
		body = "OFF"
	default:
		return fmt.Errorf("Unsupported SIA command for pusher (%s)\n", sia.command)
	}
	itemUrl := strings.Join([]string{"alarm/zone_", sia.zone, "/state"}, "")
	client.Publish(itemUrl, 0, false, body)
	
	log.Printf("Publish %s to %s", body, itemUrl)
	return nil
}