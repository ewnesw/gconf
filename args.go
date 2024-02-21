package main

import (
	"os"
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func argParser(w *fsnotify.Watcher,arg []string){
	switch arg[0]{
		case "help":
			larg := len(arg[1:])
			if larg==0{
				help("",0)
			}else if larg==1 && arg[1]=="add"{
				help("add",0)
			}else{
				help("",1)
			}
		case "add":
			if len(arg[1:])>1{
				for i := 0;  i< len(arg[1:]); i++ {
    					addPath(w,arg[i+1])
    				}
    			}else{
				addPath(w,arg[1])	
			}
		case "init":
			checkSetUp()	
		default:
			fmt.Println("unmatching argument")
			help("",1)
	}
}

func argHandler(w *fsnotify.Watcher) {
	if len(os.Args[1:]) == 0{
		help("",0)
	}else{
		argParser(w,os.Args[1:])
	}
}
