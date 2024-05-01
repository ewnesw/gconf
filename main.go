package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

func addPath(w *fsnotify.Watcher, path string) {
	err := w.Add(path)
	if err!=nil{
		fmt.Println(err)
		return
	}
	log.Println(w.WatchList())
	copyDir(path,w)
}

func watchLoop(w *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Has(fsnotify.Write) {
				log.Println("modified file:", event.Name)
				dir,filename,err := getDirFile(event.Name)
				if err != nil{
					fmt.Println(err)
					return
				}
				fmt.Println(dir+"/"+filename)
				copyfile(event.Name)
			} else if event.Has(fsnotify.Create) {
				fileinfo, err := os.Stat(event.Name)
				if err!=nil{
					fmt.Println(err)
					return
				}
				dir := fileinfo.IsDir()
				if dir {
					addPath(w, event.Name)
				}else{
					copyfile(event.Name)
				}
				log.Println("created file:", event.Name, dir)
			}
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	errorCheckFatal(err)

	defer watcher.Close()

	go watchLoop(watcher)

	argHandler(watcher)
	<-make(chan struct{})
}
