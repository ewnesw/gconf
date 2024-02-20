package main
import (
	"fmt"
	"strings"
	"log"
	"os"
)

func errorCheck(err error){
	if err!=nil{
		fmt.Println(err)
		return
	}
}

func errorCheckFatal(err error){
	if err != nil {
 		log.Fatal(err)
 	}
}

func formatPath(path string) []string{
	return strings.Split(path,"/")
}

func help(cmd string,status int){
	switch cmd {
		case "add":
			fmt.Println("example: gconf add <path>")
			os.Exit(status)
 		default:
			fmt.Println("this is the help page gl hf")
			os.Exit(status)
 	}
}
