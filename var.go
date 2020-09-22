package main

import (
	"os"
)

var (
	encoders       = []string{"ffmpeg", "null"}
	extensions     = []string{".flv", ".ts", ".m3u8"}
	workingPath, _ = os.Getwd()
	sub            = false
)
