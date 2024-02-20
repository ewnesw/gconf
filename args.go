package main

import (
	"os"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func help(status int){
	fmt.Println("example: gconf add <path>")
	os.Exit(status)
}

func argParser(w *fsnotify.Watcher,arg []string){
	switch arg[0]{
		case "help":
			if len(arg)>1{
				help(1)
			}else{
				help(0)
			}
		case "add":
			if len(arg[1:])>1{
				for i := 0;  i< len(arg[1:]); i++ {
    					addPath(w,arg[i+1])
    				}
    			}else{
				addPath(w,arg[1])	
			}
		default:
			fmt.Println("unmatching argument")
			help(1)
	}
}

func argHandler(w *fsnotify.Watcher) {
	if len(os.Args[1:]) == 0{
		help(0)
	}else{
		argParser(w,os.Args[1:])
	}
}
