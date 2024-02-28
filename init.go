package main

import (
	"fmt"
	"os"
)

func askCreation() bool {
	var rep rune
	fmt.Println("gconf about to create .gconf/backup in your home directory, continue ? [Y/n]")
	for {
		_, err := fmt.Scanf("%c", &rep)
		errorCheckFatal(err)
		if rep == rune('y') {
			return true
		} else if rep == rune('n') {
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
