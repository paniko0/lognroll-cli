package main

import (
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			ev := <-watcher.Event
			log.Println(ev)
			if ev.IsModify() {
				bb, err := ioutil.ReadFile(ev.Name)
				if nil != err {
					return
				}
				log.Println("new contents:", string(bb))
			}
		}
	}()

	err = watcher.Watch("/Users/danillosouza/dev/go/src/github.com/paniko0/lognroll-cli/teste.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	watcher.Close()
}
