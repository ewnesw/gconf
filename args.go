package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
)

func argParser(w *fsnotify.Watcher, arg []string) {
	switch arg[0] {
	case "help":
		larg := len(arg[1:])
		if larg == 0 {
			help("", 0)
		} else if larg == 1 && arg[1] == "add" {
			help("add", 0)
		} else {
			help("", 1)
		}

	case "add":
		checkSetUp()
		if len(arg[1:]) > 1 {
			for i := 0; i < len(arg[1:]); i++ {
				addPath(w, formatPath(arg[i+1]))
			}
		} else {
			addPath(w, formatPath(arg[1]))
		}

	case "init":
		checkSetUp()
		os.Exit(0)

	default:
		fmt.Println("unmatching argument")
		help("", 1)
	}
}

func argHandler(w *fsnotify.Watcher) {
	if len(os.Args[1:]) == 0 {
		help("", 0)
	} else {
		argParser(w, os.Args[1:])
	}
}
