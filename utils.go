package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
)

func errorCheckFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func formatPath(path string) (dirpath string,filename string) {
	temp := strings.Split(path, "/")
	return path[:len(path)-len(temp[len(temp)-1])], temp[len(temp)-1]
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

func checkDir(user *user.User, dir string) bool {
	_, err := os.Stat(user.HomeDir + "/" + dir)
	if os.IsNotExist(err) {
		fmt.Println("missing  " + dir + " directory")
		return false
	}
	return true
}

func createDir(user *user.User, path string) {
	err := os.MkdirAll(user.HomeDir+"/"+path, 0755)
	errorCheckFatal(err)
	fmt.Println("created " + user.HomeDir + path)
}
