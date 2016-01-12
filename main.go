package main

import (
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/pubnub/go/messaging"
	"log"
	"time"
)

var (
	logFile      string
	pubKey       string
	subscribeKey string
	secretKey    string
	channelName  string
)

func main() {

	initVars()

	done := make(chan bool)

	// Process events
	go func() {

		t, err := tail.TailFile(logFile, tail.Config{Follow: true})
		for line := range t.Lines {
			log.Println(line.Text)

			pubInstance := messaging.NewPubnub(pubKey, subscribeKey, secretKey, "", false, channelName)

			var errorChannel = make(chan []byte)
			var callbackChannel = make(chan []byte)
			go pubInstance.Publish(channelName, line.Text, callbackChannel, errorChannel)
			go handleResult(callbackChannel, errorChannel, 1000, "Publish")
			// please goto the top of this file see the implementation of handleResult

		}

		if nil != err {
			return
		}
	}()

	// Hang so program doesn't exit
	<-done

}

func handleResult(successChannel, errorChannel chan []byte, timeoutVal int64, action string) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(timeoutVal) * time.Second)
		timeout <- true
	}()
	for {
		select {
		case success, ok := <-successChannel:
			if !ok {
				break
			}
			if string(success) != "[]" {
				log.Println(fmt.Sprintf("%s Response: %s ", action, success))
				log.Println("")
			}
			return
		case failure, ok := <-errorChannel:
			if !ok {
				break
			}
			if string(failure) != "[]" {
				// if displayError {
				// 	log.Println(fmt.Sprintf("%s Error Callback: %s", action, failure))
				// 	log.Println("")
				// }
			}
			return
		case <-timeout:
			log.Println(fmt.Sprintf("%s Handler timeout after %d secs", action, timeoutVal))
			log.Println("")
			return
		}
	}
}

func initVars() {
	flag.StringVar(&pubKey, "pubKey", "", "pubnub public key")
	flag.StringVar(&subscribeKey, "subscribeKey", "", "pubnub subscribe key")
	flag.StringVar(&secretKey, "secretKey", "", "pubnub secret key")
	flag.StringVar(&logFile, "logFile", "/var/logs/error.log", "complete path to log")
	flag.StringVar(&channelName, "channelName", "", "channel name to connect to")

	flag.Parse()
}
