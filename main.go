package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Config struct {
	FilesLimit int
	LoadDir    string
}

func main() {
	config := Config{FilesLimit: 21, LoadDir: "."}
	fmt.Println("config=", config)

	current, err := filepath.Abs(config.LoadDir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//debug
	fmt.Println("current absolute path = ", current)

	files, err := os.ReadDir(current)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var targetFiles []os.FileInfo
	for _, entry := range files {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}

		targetFiles = append(targetFiles, info)
	}

	if len(targetFiles) > 21 {
		sort.Slice(targetFiles, func(i, j int) bool { return targetFiles[i].ModTime().Before(targetFiles[j].ModTime()) })
		//debug
		for index, info := range targetFiles[21:] {
			fmt.Println("index=", index, "name=", info.Name(), "info(ModTime)=", info.ModTime())
		}
	} else {
		fmt.Println("Size of the specified dir files =", len(targetFiles))
	}
}
