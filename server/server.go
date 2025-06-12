package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	fmt.Printf("Server listening at: %s:%d\n", ip, port)
	return server
}

func (server *Server) Handler(conn net.Conn) {	
	user := NewUser(conn)
	user.Online(server)
	
	isLive := make(chan bool)
	
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)

			if err != nil && err != io.EOF{
				fmt.Println("Conn read err:", err)
				return
			}

			// Read message from user
			msg := string(buf[:n-1])
			if msg != "" {
				user.DoMessage(msg, server)
			}
			isLive <- true
		}
	}()

	for {
		select {
			case <-isLive:
				// Do nothing to reset timer
			case <-time.After(5 * 60 * time.Second):
				// Force offline current user
				user.SendMessage("You are offline due to timeout")
				time.Sleep(1 * time.Second)
				user.Offline(server)
				close(user.C)
				user.conn.Close()
				return
		}
	}
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

func (server *Server) ListOnlineUsers() []User {
	var onlineUsers []User
	server.mapLock.RLock()
	for _, user := range server.OnlineMap {
		onlineUsers = append(onlineUsers, *user)
	}
	server.mapLock.RUnlock()
	return onlineUsers
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