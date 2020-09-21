package main

import (
	"container/list"
	"flag"
	"fmt"
	"github.com/cekys/gopkg"
	"os"
	"os/exec"
	"strings"
)

// setupFlags 设置命令行参数处理器
func setupFlags() {
	// 设置p 与 sub参数
	flag.StringVar(&workingPath, "p", "./", "working path")
	flag.BoolVar(&sub, "sub", false, "sub folder")
	flag.Parse()
}

// checkEncoder 检查编码器是否存在
func checkEncoder() error {
	// 编码器不存在个数的计数器
	var c = 0
	for n, i := range encoders {
		path, err := gopkg.ProgramExist(i)
		if err != nil {
			c++
			continue
		}
		// 若编码器存在则将其地址替换为绝对路径
		encoders[n] = path
	}
	// 编码器不存在个数的计数器 = 编码器设置的个数 则返回错误
	if c == len(encoders) {
		return fmt.Errorf("can't find encoders")
	}
	return nil
}

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

	for i := 0; i < fileList.Len(); i++ {
		inputFile := listNode.Value.(string)
		outputFile := gopkg.StringTrimSuffix(inputFile) + ".mp4"

		cmd := exec.Command(encoders[0], "-i", inputFile, "-codec", "copy", outputFile)
		err := gopkg.ProgramRealtimeOutput(cmd)
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
		listNode = listNode.Next()
	}
	return nil
}
