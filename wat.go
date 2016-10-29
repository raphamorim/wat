package main

import (
	"fmt"
	"gopkg.in/fsnotify.v1"
	"io/ioutil"
	"log"
	"os/exec"
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

	fmt.Println("Waiting changes...")
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// fmt.Println(event.Op)
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					go func() {
						cmd := exec.Command(w.exec)
						out, err := cmd.StdoutPipe()
						if err := cmd.Start(); err != nil {
							log.Fatal(err)
						}
						if err != nil {
							log.Fatal(err)
						}
						grepBytes, _ := ioutil.ReadAll(out)
						if err := cmd.Wait(); err != nil {
							log.Fatal(err)
						}
						fmt.Printf("%s\n", string(grepBytes))
					}()
				}
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
