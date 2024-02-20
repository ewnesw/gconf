package main
import (
	"fmt"
	"strings"
	"log"
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
