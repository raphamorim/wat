package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) <= 2 {
		os.Exit(2)
	}

	w, err := newWatch(os.Args[1], os.Args[2], os.Args[3:], os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	defer w.close()

	// Wait for CTRL-C
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
