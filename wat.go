package main

import (
	"bufio"
	"fmt"
	"gopkg.in/fsnotify.v1"
	"io"
	"log"
	"os"
	"os/exec"
)

func startWatch(path, command string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	fmt.Println("Waiting...")
	done := make(chan bool)
	go func(watcher *fsnotify.Watcher, command string) {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					go func() {
						cmd := exec.Command(command)
						var out io.ReadCloser
						out, err = cmd.StdoutPipe()
						if err != nil {
							log.Fatal(err)
						}

						scanner := bufio.NewScanner(out)
						go func() {
							for scanner.Scan() {
								fmt.Printf("%s\n", scanner.Text())
							}
						}()
						if err = cmd.Run(); err != nil {
							log.Fatal(err)
							os.Exit(1)
						}
					}()
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println(event.Name)
				}
			case err = <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}(watcher, command)

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
