package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C  chan string
	conn net.Conn
}

// Create a user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	
	user := &User {
		Name: userAddr,
		Addr: userAddr,
		C: make(chan string),
		conn: conn,
	}
	
	go user.ListenForMessage()
	
	return user
}

func (user *User) Online(server *Server) {
	// Add current user to onlineMap
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()
	// Brodcast user online Message
	server.Brodcast(user, "online")
}

func (user *User) Offline(server *Server) {
	// Remove current user from onlineMap
	server.mapLock.Lock()
	delete(server.OnlineMap, user.Name)
	server.mapLock.Unlock()
	// Brodcast user offline Message
	server.Brodcast(user, "offline")
}

func (user *User) SendMessage(msg string, server *Server) {
	switch msg {
		case "/list":
			onlineUsers := server.ListOnlineUsers()
			var builder strings.Builder
			for i, user := range onlineUsers {
				builder.WriteString(user.Name)
				if i < len(onlineUsers)-1 {
					builder.WriteString("\n")
				}
			} 
			user.C <- builder.String()
		default: server.Brodcast(user, msg)
	}
}

// lisen current user's channel
func (user *User) ListenForMessage() {
	for {
		msg := <-user.C
		
		user.conn.Write([]byte(msg + "\n"))
	}
}