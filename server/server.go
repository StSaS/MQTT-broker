package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	//"math/rand"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"sync"
)

var (
	in        string
	counter   int
	words     []string
	startTime time.Time
	duration  time.Duration
	Mu        = &sync.Mutex{}
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	/*fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())*/
	Mu.Lock()
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
	Mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			Mu.Unlock()
			fmt.Println(r)
		}
	}()

}

func Receive(word string, c MQTT.Client) {

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
	opts := MQTT.NewClientOptions().AddBroker("tcp://mosquitto:1883")
	//opts := MQTT.NewClientOptions().AddBroker("tcp://iot.eclipse.org:1883")
	opts.SetClientID("server")
	opts.SetDefaultPublishHandler(f)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, s := range words {
		go Receive(s, c)
	}

	startTime = time.Now()

	fmt.Scanln()
}
