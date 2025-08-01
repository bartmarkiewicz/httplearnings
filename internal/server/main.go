package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

const port = 42069

type Server struct {
	port          int
	serverRunning atomic.Bool
	listener      net.Listener
}

func (s *Server) Serve(port int) (*Server, error) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s.listener = listen
	s.serverRunning.Store(true)
	err = s.listen()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Close() error {
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(s.listener)
	s.serverRunning.Store(false)
	return nil
}

func (s *Server) handle(conn net.Conn) {
	_, err := conn.Write([]byte("HTTP/1.1 200 OK\nContent-Type: text/plain\n\nHello World!"))
	if err != nil {
		return
	}
	println("Connection accepted")
	defer conn.Close()
	println("Connection closed")
}

func (s *Server) listen() error {
	for {
		if !s.serverRunning.Load() {
			return nil
		}
		acceptedConnection, err := s.listener.Accept()
		go s.handle(acceptedConnection)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
}

func main() {
	server := &Server{}

	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
