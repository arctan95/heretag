package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int
	OnlineMap map[string]*User
	mapLock sync.RWMutex
	Message chan string
}

// Create a server
func NewServer(ip string, port int) *Server {
	server := &Server {
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	fmt.Printf("Server listening at %s:%d\n", ip, port)
	return server
}

func (server *Server) Handler(conn net.Conn) {
	
	user := NewUser(conn)
	
	// Add current user to onlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()
	
	// Brodcast user online Message
	server.Brodcast(user, "online")
}

func (server *Server) ListenServerMessgae() {
	for {
		msg := <-server.Message
		
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Brodcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

// Start server
func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println("Failed to listen:", err)
	}
	
	defer listener.Close()
	
	go server.ListenServerMessgae()
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener accept err:", err)
			continue
		}
		
		go server.Handler(conn)
	}
}