package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

// Create a server
func NewServer(ip string, port int) *Server {
	server := &Server {
		Ip: ip,
		Port: port,
	}
	fmt.Printf("Server listening at %s:%d\n", ip, port)
	return server
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Println("Connection established successfully")
}

// Start server
func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("Failed to listen:", err)
	}
	
	defer listener.Close()
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		
		go server.Handler(conn)
	}
}