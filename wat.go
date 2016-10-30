package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"

	"gopkg.in/fsnotify.v1"
)

type watch struct {
	watcher *fsnotify.Watcher
}

func (w *watch) close() error {
	return w.watcher.Close()
}

func newWatch(path, command string, args []string, stdout io.Writer) (*watch, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(stdout, "Waiting...")
	go func(watcher *fsnotify.Watcher, command string) {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					go func() {
						cmd := exec.Command(command, args...)
						cmd.Stdout = stdout
						if err = cmd.Run(); err != nil {
							log.Println(err)
						}
					}()
				}
			case err = <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}(watcher, command)

	err = watcher.Add(path)
	if err != nil {
		return nil, err
	}
	return &watch{watcher}, nil
}
