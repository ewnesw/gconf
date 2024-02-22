package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"io"
	"sync"
	"strings"
)

func pushUpdate(directory string, commit_msg string) {
	fmt.Println(directory)
	r, err := git.PlainOpen(directory)
	if err!=nil{
		return
	}
	w, err := r.Worktree()
	if err!=nil{
		return
	}
	_, err = w.Add(".")
	if err!=nil{
		return
	}
	status, err := w.Status()
	if err!=nil{
		return
	}
	fmt.Println(status)
	_, err = w.Commit(commit_msg, &git.CommitOptions{})
	if err!=nil{
		return
	}
	err = r.Push(&git.PushOptions{})
	if err!=nil{
		return
	}

}

func copyfile(filepath string){
	temp := strings.Split(filepath,"/")
	filename := temp[len(temp)-1]
	fmt.Println("copying " + filename)
	user := getUser()
	src,err := os.Open(filepath)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer  src.Close()

	dst,err :=  os.Create(user.HomeDir+"/.gconf/backup/"+filename)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer  dst.Close()

	 _, err = io.Copy(dst, src) 
	if err!=nil{
		fmt.Println(err)
		return
	}

    	err = dst.Sync()
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println("done")
}

func lookupDir(path string) ([]string,error){
	var filelist []string
	f, err := os.Open(path)
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
    	files, err := f.Readdir(0)
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
    	for _,v := range files {
		if !v.IsDir(){
			filelist = append(filelist,v.Name())
		}
    	}
	fmt.Println(filelist)
	return filelist,nil
}

func copyDir(path string){
	var wg sync.WaitGroup
	fl, err := lookupDir(path)
	if err!=nil{
		fmt.Println(err)
		return
	}
	for _,v := range fl {
		wg.Add(1)
		go func(filepath string){
			defer wg.Done()
			copyfile(filepath)
		}(path+"/"+v)
	}
	wg.Wait()
}


func addPath(w *fsnotify.Watcher, path string) {
	err := w.Add(path)
	if err!=nil{
		fmt.Println(err)
		return
	}
	log.Println(w.WatchList())
	copyDir(path)
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
				paths := formatPath(event.Name)
				fmt.Println(paths)
			} else if event.Has(fsnotify.Create) {
				fileinfo, err := os.Stat(event.Name)
				if err!=nil{
					fmt.Println(err)
					return
				}
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
