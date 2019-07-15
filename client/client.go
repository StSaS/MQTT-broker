package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

import MQTT "github.com/eclipse/paho.mqtt.golang"

func Publish(word string) {

	opts := MQTT.NewClientOptions().AddBroker("tcp://iot.eclipse.org:1883")
	opts.SetClientID("client_topic_" + word)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	topic := "topic_" + word
	rand.Seed(time.Now().UnixNano())

	for {

		token := c.Publish(topic, 0, false, word)
		token.Wait()
		n := rand.Intn(10) // n will be between 0 and 10
		fmt.Print("Publish to " + topic + "\t")
		fmt.Printf("Sleeping %d seconds...\n", n)
		time.Sleep(time.Duration(n) * time.Second)

	}

	//defer c.Disconnect(250)
}

func main() {

	in := os.Getenv("in")
	fmt.Println("varible(in):", in)

	if in == "" {
		fmt.Println("Please, set varible in. On default in=Hello world")
		in = "Hello world"
	}

	words := strings.Split(in, " ")

	for _, s := range words {
		go Publish(s)
	}

	fmt.Scanln()
}
