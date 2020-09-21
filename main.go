package main

import (
	"container/list"
	"github.com/cekys/gopkg"
	"log"
	"path/filepath"
)

func main() {

	setupFlags()

	err := checkEncoder()
	if err != nil {
		log.Fatal(err)
	}

	fileList := list.New()
	inputList := list.New()
	var filter []string

	//Find all files in working path
	err = gopkg.PathWalk(workingPath, fileList, filter, sub)
	if err != nil {
		log.Fatal(err)
	}

	//Find all files in the directory with the specified suffix
	inputList = findFilesWithSuffix(fileList, extensions)

	//Check if input files exists
	if inputList.Len() == 0 {
		p, err := filepath.Abs(workingPath)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalln("No files to convert in", p)
	}

	//Convert files and delete input files when conversion is successful
	err = convert2mp4(inputList, 1)
	if err != nil {
		log.Fatal(err)
	}
}
