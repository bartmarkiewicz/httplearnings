package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {

	listen, err := net.Listen("tcp", ":42069")
	defer listen.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		accepted, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Connection accepted")
		linesCh := getLinesChannel(accepted)

		for line := range linesCh {
			fmt.Println(line)
		}

		fmt.Println("Connection closed")
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					lines <- currentLineContents
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
			str := string(b[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
				currentLineContents = ""
			}
			currentLineContents += parts[len(parts)-1]
		}
	}()
	return lines
}
