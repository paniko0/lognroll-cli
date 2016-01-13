package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pubnub/go/messaging"
	"io"
	"log"
	"os"
	"time"
)

var (
	pubKey       string
	subscribeKey string
	secretKey    string
	channelName  string
)

func main() {

	initFlags()

	done := make(chan bool)

	// Process events
	go func() {

		in := bufio.NewReader(os.Stdin)
		stats, err := os.Stdin.Stat()

		if err != nil {
			fmt.Println("file.Stat()", err)
		}

		if stats.Size() > 0 {

			for {
				input, err := in.ReadString('\n')
				if err != nil && err == io.EOF {
					break
				}

				pubInstance := messaging.NewPubnub(pubKey, subscribeKey, secretKey, "", false, channelName)

				var errorChannel = make(chan []byte)
				var callbackChannel = make(chan []byte)
				go pubInstance.Publish(channelName, input, callbackChannel, errorChannel)
				go handleResult(callbackChannel, errorChannel, 1000, "Publish")

				if nil != err {
					return
				}

			}

		}

	}()

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

func initFlags() {
	flag.StringVar(&pubKey, "pubKey", "", "pubnub public key")
	flag.StringVar(&subscribeKey, "subscribeKey", "", "pubnub subscribe key")
	flag.StringVar(&secretKey, "secretKey", "", "pubnub secret key")
	flag.StringVar(&channelName, "channelName", "", "channel name to connect to")

	flag.Parse()
}
