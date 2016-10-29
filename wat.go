package main

import (
	"bufio"
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"
	"os"
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

	fmt.Println("Waiting...")
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
						if err != nil {
							log.Fatal(err)
						}

						scanner := bufio.NewScanner(out)
						go func() {
							for scanner.Scan() {
								fmt.Printf("%s\n", scanner.Text())
							}
						}()
						if err := cmd.Start(); err != nil {
							log.Fatal(err)
							os.Exit(1)
						}
						if err := cmd.Wait(); err != nil {
							log.Fatal(err)
							os.Exit(1)
						}
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
