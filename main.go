package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	open, err := os.Open("messages.txt")
	if err != nil {
		return
	}

	defer open.Close()

	for each := range getLinesChannel(open) {
		fmt.Printf("read: %s\n", each)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	currentLine := ""

	go func() {
		bytes := make([]byte, 8)
		defer f.Close()
		r, err := f.Read(bytes)

		for {
			currentLine += string(bytes)

			if errors.Is(err, io.EOF) || r == 0 {
				break
			}

			lines := strings.Split(currentLine, "\n")

			if len(lines) > 1 {
				for _, e := range lines[0 : len(lines)-1] {
					ch <- e
					currentLine = lines[len(lines)-1]
				}
			}
			_, err = f.Read(bytes)
		}
		if currentLine != "" {
			fmt.Println(currentLine)
			ch <- currentLine
		}
		close(ch)

	}()

	return ch
}
