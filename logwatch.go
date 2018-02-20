package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"
	"time"
)

// Just for fun, not used: print a whole text file using bufio.Scanner
func PrintFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type LogEvent struct {
	source string
	event  string
}

// WatchFile watches a file like "tail -f" does. It cannot deal with cases
// such as the file being moved or deleted.
func WatchFile(filename string, pattern string, ch chan LogEvent, wg *sync.WaitGroup) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal("Invalid regexp: ", err)
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Seek(0, 2) // Start at current end of file
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		switch err {
		case nil:
			if re.MatchString(line) {
				ch <- LogEvent{source: filename, event: string(line)}
			}
		case io.EOF:
			time.Sleep(1 * time.Second)
		default: // other error
			break
		}
	}
	wg.Done()
}

func consume(ch chan LogEvent) {
Loop:
	for {
		select {
		case line, ok := <-ch:
			if !ok {
				break Loop
			}
			fmt.Printf("%s %s", line.source, line.event)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: logwatch <file1> <pattern1> [<file2> <pattern2> ...]")
		os.Exit(1)
	}
	ch := make(chan LogEvent)
	var wg sync.WaitGroup
	for i := 1; i < len(os.Args); i += 2 {
		wg.Add(1)
		go WatchFile(os.Args[i], os.Args[i+1], ch, &wg)
	}
	go consume(ch)
	wg.Wait()
	close(ch)
}
