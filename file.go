package main

import (
	"fmt"
	"os"
	"github.com/fsnotify/fsnotify"
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

func lookupDir(path string, w *fsnotify.Watcher) ([]string,error){
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
		}else{
			fmt.Println("lookup " + path + "/" +v.Name())
			addPath(w,path + "/" +v.Name())
		}
    	}
	fmt.Println(filelist)
	return filelist,nil
}


func copyDir(path string, w *fsnotify.Watcher){
	var wg sync.WaitGroup
	fl, err := lookupDir(path,w)
	if err!=nil{
		fmt.Println(err)
		return
	}
	_, filename,err:=getDirFile(path) // we need the name of the path given that's why wee treat it as a file not a directory
	errorCheckFatal(err)
	fmt.Println("path: " + path + " dir: " + filename)
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
