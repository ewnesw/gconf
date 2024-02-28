package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func askCreation() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("gconf about to create .gconf/backup in your home directory, continue ? [Y/n]")
	for {
		rep, err := reader.ReadString('\n')
		errorCheckFatal(err)
		rep = strings.TrimSuffix(rep, "\n")
		if rep == "" || rep == "y" {
			return true
		} else if rep == "n" {
			return false
		} else {
			fmt.Println("pls enter y or n")
		}
	}
}

func checkSetUp() {
	user := getUser()
	if !checkDir(user.HomeDir+"/.gconf/backup") {
		if askCreation() {
			createDir(user.HomeDir+"/.gconf/backup")
		} else {
			fmt.Println("Aborting")
			os.Exit(0)
		}
	} else {
		fmt.Println("Everything is fine here")
	}
}
