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

func (user *User) SendMessage(msg string) {
	user.C <- msg
}

func (user *User) DoMessage(msg string, server *Server) {
	switch {
		case msg == "/list":
			onlineUsers := server.ListOnlineUsers()
			var builder strings.Builder
			for i, user := range onlineUsers {
				builder.WriteString(user.Name)
				if i < len(onlineUsers)-1 {
					builder.WriteString("\n")
				}
			} 
			user.SendMessage(builder.String())
		case strings.Contains(msg, "/rename"):
			newName := strings.Split(msg, " ")[1]
			// Check if username has been taken
			_, ok := server.OnlineMap[newName]
			if (ok) {
				user.SendMessage("This username has been taken by someone!")
			} else {
				server.mapLock.Lock()
				delete(server.OnlineMap, user.Name)
				server.OnlineMap[newName] = user
				server.mapLock.Unlock()
				
				user.Name = newName
				user.SendMessage("Your new username is: " + user.Name)
			}
			
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