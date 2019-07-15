package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	//"math/rand"
)

import MQTT "github.com/eclipse/paho.mqtt.golang"

var (
	in        string
	counter   int
	words     []string
	startTime time.Time
	duration  time.Duration
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	/*fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())*/
	if (counter == len(words)) && (words[counter-1] == string(msg.Payload())) {
		duration = time.Since(startTime)
		fmt.Println(duration.String() + ": " + in)
		startTime = time.Now()
		counter = 1
	} else {
		if words[counter-1] == string(msg.Payload()) {
			counter += 1
		} else {
			startTime = time.Now() //incorrect sequence
		}
	}

}

func Receive(word string) {

	opts := MQTT.NewClientOptions().AddBroker("tcp://iot.eclipse.org:1883")
	opts.SetClientID("server_topic_" + word)
	opts.SetDefaultPublishHandler(f)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "topic_" + word
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)

	}

	fmt.Println(topic)

	for {
		fmt.Print()
	}

}

func main() {

	in = os.Getenv("in")
	fmt.Println("varible(in):", in)

	if in == "" {
		fmt.Println("Please, set varible in. On default in=Hello world")
		in = "Hello world"
	}

	words = strings.Split(in, " ")
	counter = 1

	for _, s := range words {
		go Receive(s)
	}

	startTime = time.Now()

	fmt.Scanln()
}