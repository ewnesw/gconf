package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/user"
	"io"
	"sync"
)

func checkPerm(user *user.User, dirfile string, filepath string){
	srcinfo, err := os.Stat(filepath)
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(srcinfo.Mode())
	err = os.Chmod(user.HomeDir+"/.gconf/backup"+dirfile,srcinfo.Mode())
	if err!=nil{
		fmt.Println(err)
		return
	}
}

func copyfile(filepath string) {
	dir, filename, err := getDirFile(filepath)
	if err!=nil{
		fmt.Println(err)
		return
	}
	fmt.Println(filepath + "copying " + dir + "/" + filename)
	user := getUser()
	src,err := os.Open(filepath)
	if err!=nil{
		fmt.Println(err)
		return
	}
	defer  src.Close()

	dst,err :=  os.Create(user.HomeDir+"/.gconf/backup"+dir + "/" + filename)
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
	fmt.Println("done copying checking perm")
	checkPerm(user,dir + "/" + filename,filepath)
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
	_, filename,err:=getDirFile(path) // we need the name of the path given that's why wee treat it as a file not a directory
	if err!=nil{	                  // kinda dogshit need to work on it and var/func names
		fmt.Println(err)
		return
	}
	fmt.Println("path: " + path + "dir: " + filename)
	createDir(getUser().HomeDir +"/.gconf/backup/"+filename)
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
