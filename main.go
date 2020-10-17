package main

import (
	"container/list"
	"flag"
	"github.com/cekys/gopkg"
	"log"
	"os"
	"path/filepath"
)

var (
	encoders       = []string{"ffmpeg"}
	extensions     = []string{".flv", ".ts", ".m3u8"}
	workingPath, _ = os.Getwd()
	sub            = false
)

func main() {
	// 设置flag
	flag.StringVar(&workingPath, "p", "./", "working path")
	flag.BoolVar(&sub, "sub", false, "sub folder")
	flag.Parse()

	// 检查编码器是否存在
	for i, encoder := range encoders {
		encoder, err := gopkg.ProgramExist(encoder)
		if err != nil {
			log.Fatal(err)
		}
		encoders[i] = encoder
	}

	fileList := list.New()
	filter := make([]string, 0)

	//Find all files in working path
	err := gopkg.PathWalk(workingPath, fileList, filter, sub)
	if err != nil {
		log.Fatal(err)
	}

	//Find all files in the directory with the specified suffix
	inputList := findFilesWithSuffix(fileList, extensions)

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
