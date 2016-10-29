package main

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"
)

type wat struct {
	path, exec string
}

func (w *wat) getCommand() string {
	return w.exec
}

func (w *wat) getPath() string {
	return w.path
}

func (w *wat) startWatch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println(event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(w.path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
