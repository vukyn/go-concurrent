package concurrent

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/XANi/loremipsum"
)

func CountWordInFile(de fs.DirEntry, path string, word string) int {
	if de.IsDir() {
		// Scan all files and folders in directory
		folderPath := fmt.Sprintf("%v/%s", path, de.Name())
		dir, err := os.ReadDir(folderPath)
		if err != nil {
			log.Fatal(err)
		}
		wg := &sync.WaitGroup{}
		count := 0
		countCh := make(chan int, len(dir))
		for _, e := range dir {
			if de.IsDir() {
				countCh <- CountWordInFile(e, folderPath, word)
			} else {
				wg.Add(1)
				go func(e fs.DirEntry) {
					defer wg.Done()
					countCh <- countWord(fmt.Sprintf("%v/%s", folderPath, e.Name()), word)
				}(e)
			}
		}
		wg.Wait()
		close(countCh)
		for c := range countCh {
			count += c
		}
		return count
	} else {
		return countWord(fmt.Sprintf("%v/%s", path, de.Name()), word)
	}
}

func countWord(filePath string, word string) int {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Cannot open file")
	}
	defer file.Close()

	// Scan file word by word
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		// Remove all special characters and convert to lower case
		if regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(strings.ToLower(scanner.Text()), "") == word {
			count++
		}
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
	return count
}

func GenFiles(weight int, path string) {
	loremIpsumGenerator := loremipsum.New()
	for i := 1; i <= weight; i++ {
		paragraph := loremIpsumGenerator.Paragraph()
		if err := WriteFile(paragraph, fmt.Sprintf("%v/Text%d.txt", path, i)); err != nil {
			log.Fatal(err)
		}
	}
}

func GenFilesAndFolders(weight int, path string) {
	var filePath string
	loremIpsumGenerator := loremipsum.New()
	for i := 1; i <= weight; i++ {
		paragraph := loremIpsumGenerator.Paragraph()
		if i%3 == 0 {
			filePath = fmt.Sprintf("%s/%s", path, loremIpsumGenerator.Word())
			if err := WriteFile(paragraph, fmt.Sprintf("%s/Text%d.txt", filePath, i)); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if i%7 == 0 {
			filePath = fmt.Sprintf("%s/%s", path, loremIpsumGenerator.Word())
			if err := WriteFile(paragraph, fmt.Sprintf("%s/Text%d.txt", filePath, i)); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if i%13 == 0 {
			filePath = fmt.Sprintf("%s/%s", path, loremIpsumGenerator.Word())
			if i%2 == 0 {
				filePath = fmt.Sprintf("%s/%s", filePath, loremIpsumGenerator.Word())
				if err := WriteFile(paragraph, fmt.Sprintf("%s/Text%d.txt", filePath, i)); err != nil {
					log.Fatal(err)
				}
			} else if err := WriteFile(paragraph, fmt.Sprintf("%s/Text%d.txt", filePath, i)); err != nil {
				log.Fatal(err)
			}
			continue
		}
		if err := WriteFile(paragraph, fmt.Sprintf("%v/Text%d.txt", path, i)); err != nil {
			log.Fatal(err)
		}
	}
}
