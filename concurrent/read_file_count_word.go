package concurrent

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	FOLDER_PATH = "assets/read_file_count_word"
	DETECT_WORD = "ut"
)

func ReadFileCountWord() {
	num := 3
	wg := &sync.WaitGroup{}
	countCh := make(chan int, 1)
	countCh <- 0
	wg.Add(num)
	for i := 1; i <= num; i++ {
		go func(i int) {
			// Open file
			file, err := os.Open(fmt.Sprintf("%v/Text%d.txt", FOLDER_PATH, i))
			if err != nil {
				log.Fatalf("Cannot open file")
			}
			defer file.Close()

			// Scan file word by word
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanWords)
			for scanner.Scan() {
				if strings.ToLower(scanner.Text()) == DETECT_WORD {
					count := <-countCh
					countCh <- count + 1
				}
			}
			if scanner.Err() != nil {
				log.Fatal(scanner.Err())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	log.Printf("Count: %d", <-countCh)
}