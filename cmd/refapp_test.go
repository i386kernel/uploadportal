package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestIndexfile(t *testing.T) {
	getwd, _ := os.Getwd()
	fmt.Println(getwd)
	fi, err := os.Stat(getwd + "/" + "index.html")
	if err != nil {
		fmt.Println(err)
	}
	if fi.Name() != "index.html" {
		log.Fatal("index.html does not exit")
	}
}

//
//func TestUploadedFile(t *testing.T) {
//	getwd, _ := os.Getwd()
//	locFi, err := os.Stat(getwd + "/" + "uploads")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(locFi.Name())
//	fmt.Println()
//	if locFi.Name() != uploadfilepath {
//		log.Fatal("Local file does not exist")
//	}
//}

func TestFileInObjStore(t *testing.T) {
	objs3 := minios3object
	objs3.GetObjects()
}
