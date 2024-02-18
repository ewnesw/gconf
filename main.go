package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"strings"
)

func errorCheck(err error){
	if err!=nil{
		fmt.Println(err)
		return
	}
}

func errorCheckFatal(err error){
	if err!=nil{
		fmt.Println(err)
		return
	}
}

func pushUpdate(directory string,commit_msg string){
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

func formatPath(path string) []string{
	return strings.Split(path,"/")
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	errorCheckFatal(err)

	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
					paths := formatPath(event.Name)
					pushUpdate(paths[0],paths[len(paths)-1])
				} else if event.Has(fsnotify.Create) {
					fileinfo, err := os.Stat(event.Name)
					errorCheck(err)	
					dir := fileinfo.IsDir()
					if dir {
						watcher.Add(event.Name)
						errorCheck(err)	
						log.Println(watcher.WatchList())
					}
					log.Println("created file:", event.Name, dir)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("./test")
	errorCheckFatal(err)	

	<-make(chan struct{})
}
