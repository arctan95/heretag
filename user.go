package main

import "net"

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

// lisen current user's channel
func (user *User) ListenForMessage() {
	for {
		msg := <-user.C
		
		user.conn.Write([]byte(msg + "\n"))
	}
}