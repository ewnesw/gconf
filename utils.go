package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"errors"
)

func errorCheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getDirFile(path string) (dir string,filename string,err error) {
	if path[0] != '/'{
		return path,path,errors.New("path must start with '/'")
	}
	path = formatPath(path)
	temp := strings.Split(path, "/")
	if len(temp)==2{
		return "", temp[len(temp)-1],nil
	}
	return "/"+temp[len(temp)-2],temp[len(temp)-1],nil 
}

func formatPath(path string) string{
	if path[len(path)-1]=='/'{ 
		return path[:len(path)-1]
	}
	return path
}

func help(cmd string, status int) {
	switch cmd {
	case "add":
		fmt.Println("example: gconf add <path>")
		os.Exit(status)
	default:
		fmt.Println("this is the help page gl hf")
		os.Exit(status)
	}
}

func getUser() *user.User {
	user, err := user.Current()
	errorCheckFatal(err)
	return user
}

func checkDir(dirpath string) bool {
	_, err := os.Stat(dirpath)
	if os.IsNotExist(err) {
		fmt.Println("missing  " + dirpath + " directory")
		return false
	}
	return true
}

func createDir(dirpath string) {
	err := os.MkdirAll(dirpath, 0755)
	errorCheckFatal(err)
	fmt.Println("created " + dirpath)
}
