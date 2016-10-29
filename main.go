package main

import (
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		os.Exit(2)
	} else {
		startWatch(os.Args[1], os.Args[2])
	}
}
