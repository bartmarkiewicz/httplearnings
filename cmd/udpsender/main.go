package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":42069")
	if err != nil {
		fmt.Println(err)
		return
	}
	udp, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer udp.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = udp.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
