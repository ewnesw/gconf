package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
)


func pushUpdate(directory string,commit_msg string){
	fmt.Println(directory)
	r, err := git.PlainOpen(directory)
	errorCheck(err)	
	w, err := r.Worktree()
	errorCheck(err)
	_, err = w.Add(".")
	errorCheck(err)
	status, err := w.Status()
	errorCheck(err)
	fmt.Println(status)
	_, err = w.Commit(commit_msg,&git.CommitOptions{})
	errorCheck(err)
	
	err = r.Push(&git.PushOptions{})
	errorCheck(err)

}

func addPath(w *fsnotify.Watcher, path string){
	err := w.Add(path)
	errorCheck(err)	
	log.Println(w.WatchList())
}

func watchLoop(w *fsnotify.Watcher){
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return
			}
			log.Println("event:", event)
			if event.Has(fsnotify.Write) {
				log.Println("modified file:", event.Name)
				paths := formatPath(event.Name)
				fmt.Println(paths)
			} else if event.Has(fsnotify.Create) {
				fileinfo, err := os.Stat(event.Name)
				errorCheck(err)	
				dir := fileinfo.IsDir()
				if dir {
					addPath(w, event.Name)
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
