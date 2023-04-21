package main

import (
	"fmt"
	"go-concurrent/concurrent"
	"io/fs"
	"log"
	"os"
	"sync"
)

func main() {
	const (
		FOLDER_PATH = "assets/read_file_count_word"
		DETECT_WORD = "ut"
	)
	// Generate mock files and folders
	concurrent.GenFilesAndFolders(100, FOLDER_PATH)

	// Scan all files and folders in directory
	dir, err := os.ReadDir(FOLDER_PATH)
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	count := 0
	countCh := make(chan int, len(dir))

	for _, e := range dir {
		wg.Add(1)
		go func(e fs.DirEntry) {
			defer wg.Done()
			countCh <- concurrent.CountWordInFile(e, FOLDER_PATH, DETECT_WORD)
		}(e)
	}
	wg.Wait()
	close(countCh)
	for c := range countCh {
		count += c
	}
	fmt.Printf("Count [%s]: %d", DETECT_WORD, count)
}
