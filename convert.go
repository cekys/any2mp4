package main

import (
	"container/list"
	"fmt"
	"localhost/vivycore"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
)

//Find files with a specific suffix
func findFilesWithSuffix(fileList *list.List, suffix []string) *list.List {
	returnList := list.New()
	file := fileList.Front()

	for i := 0; i < fileList.Len(); i++ {
		for _, ext := range suffix {
			if strings.HasSuffix(file.Value.(string), ext) {
				returnList.PushBack(file.Value)
			}
		}
		file = file.Next()
	}
	return returnList
}

func convert2mp4(fileList *list.List, mode int) error {
	listNode := fileList.Front()
	// 按cpu数目设置最大协程数目
	max := runtime.NumCPU()
	// 协程数组
	chs := make([]chan string, max)
	// 声明协程并启动
	for i := 0; i < max; i++ {
		chs[i] = make(chan string)
		go func(ch chan string) {
			ch <- "start"
		}(chs[i])
	}

	// 声明selectCase[]
	var selectCase = make([]reflect.SelectCase, max)
	// 给每个selectCase 设置目的地和对应的信道
	for i := 0; i < max; i++ {
		selectCase[i].Dir = reflect.SelectRecv
		selectCase[i].Chan = reflect.ValueOf(chs[i])
	}

	for i := 0; i < fileList.Len(); i++ {
		inputFile := listNode.Value.(string)
		outputFile := vivycore.StringTrimSuffix(inputFile) + ".mp4"
		chosen, _, recvOk := reflect.Select(selectCase)

		if recvOk {
			go runit(inputFile, outputFile, mode, chs[chosen])
		}

		listNode = listNode.Next()
	}
	for _, ch := range chs {
		<-ch
	}
	return nil
}

func runit(inputFile string, outputFile string, mode int, ch chan string) error {
	cmd := exec.Command(encoders[0], "-i", inputFile, "-codec", "copy", "-map", "0", outputFile)
	err := vivycore.ProgramRealtimeOutput(cmd)
	if err != nil {
		return err
	}

	//If mode = 1 and the cmd completed without errors, delete the input file
	if mode == 1 {
		fmt.Println("File:", inputFile, "converted successfully")
		fmt.Println("Remove the file:", inputFile)
		err = os.Remove(inputFile)
		if err != nil {
			return err
		}
	}
	ch <- "done"
	return nil
}
