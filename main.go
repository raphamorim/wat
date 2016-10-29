package main

import (
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		os.Exit(2)
	} else {
		w := wat{
			path: os.Args[1],
			exec: os.Args[2],
		}

		w.startWatch()
	}
}
