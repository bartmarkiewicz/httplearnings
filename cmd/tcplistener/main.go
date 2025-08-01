package main

import (
	"fmt"
	"learnhttp/internal/request"
	"net"
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
		request, err := request.RequestFromReader(accepted)

		fmt.Printf("Request line: \n- Method: %s\n- Target: %s\n- Version: %s\n",
			request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)

		fmt.Printf("Headers: \n")
		for k, v := range request.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Printf("Body: \n%v\n", string(request.Body))

		fmt.Println("Connection closed")
	}

}
