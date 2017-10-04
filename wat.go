package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

type RecursiveWatcher struct {
	*fsnotify.Watcher
}

func (w *RecursiveWatcher) RecursiveAdd(name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	} else if fi.IsDir() {
		err := filepath.Walk(name, func(newPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// log.Println("added", newPath)
				w.Add(newPath)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		err := w.Add(name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *RecursiveWatcher) RecursiveRemove(name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	} else if fi.IsDir() {
		err := filepath.Walk(name, func(newPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// log.Println("removed", newPath)
				w.Remove(newPath)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		w.Remove(name)
		if err != nil {
			return err
		}
	}
	return nil
}

type watch struct {
	watcher *RecursiveWatcher
}

func (w *watch) close() error {
	return w.watcher.Close()
}

func newWatch(path, command string, args []string, stdout io.Writer) (*watch, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	recurWatcher := &RecursiveWatcher{watcher}

	fmt.Fprintln(stdout, "Waiting...")
	go func(watcher *RecursiveWatcher, command string) {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println(event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					if fi, _ := os.Stat(event.Name); fi.IsDir() {
						watcher.RecursiveAdd(event.Name)
					}
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					watcher.RecursiveRemove(event.Name)
				}
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
	}(recurWatcher, command)

	err = recurWatcher.RecursiveAdd(path)
	if err != nil {
		return nil, err
	}
	return &watch{recurWatcher}, nil
}
