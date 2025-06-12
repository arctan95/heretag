package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

// Create a new client
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag: -1,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	client.conn = conn
	fmt.Printf("Connected to server at: %s:%d\n", serverIp, serverPort)

	return client
}

func (client *Client) menu() bool {
	var choice int

	fmt.Println("1.Public chat")
	fmt.Println("2.Private chat")
	fmt.Println("3.Rename")
	fmt.Println("0.Exit")
	
	fmt.Scanln(&choice)
	if choice >= 0 && choice <= 3 {
		client.flag = choice
	} else {
		fmt.Println("Invalid choice, please try again.")
		return false
	}
	return true
}

func (client *Client) SendMessage(message string) {
	if message == "" {
		fmt.Println("Message cannot be empty.")
		return
	}
	_, err := client.conn.Write([]byte(message + "\n"))
	if err != nil {
		fmt.Println("Error sending message to server:", err)
		return
	}
}

func (client *Client) UpdateName() {
	fmt.Print("Enter your name: ")
	fmt.Scanln(&client.Name)
	if client.Name == "" {
		client.Name = "Anonymous"
	}
	sendMsg := fmt.Sprintf("/rename %s", client.Name)
	client.SendMessage(sendMsg)
}

func (client *Client) ListOnlineUsers() {
	client.SendMessage("/list")
}

func (client *Client) PrivateChat() {
	var targetUser string
	var message string
	client.ListOnlineUsers()
	fmt.Println("Enter the username of the user you want to chat with, /exit to quit: ")
	fmt.Scanln(&targetUser)

	for targetUser != "/exit" {
		fmt.Println("Enter your message, /exit to quit: ")
		fmt.Scanln(&message)

		for message != "/exit" {
			sendMsg := fmt.Sprintf("/to %s %s", targetUser, message)
			if (len(message) > 0) {
				client.SendMessage(sendMsg)
			}
			message = ""
			fmt.Println("Enter your message, /exit to quit: ")
			fmt.Scanln(&message)
		}
		client.ListOnlineUsers()
		fmt.Println("Enter the username of the user you want to chat with, /exit to quit: ")
		fmt.Scanln(&targetUser)
	}
}

func (client *Client) PublicChat() {
	var message string
	fmt.Println("Enter your message, /exit to quit: ")
	fmt.Scanln(&message)

	for message != "/exit" {
		if (len(message) > 0) {
			client.SendMessage(message)
		}
		message = ""
		fmt.Println("Enter your message, /exit to quit: ")
		fmt.Scanln(&message)
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.menu() {
		switch client.flag {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.UpdateName()
		default:
			fmt.Println("Invalid option, please try again.")
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Server IP address")
	flag.IntVar(&serverPort, "port", 6969, "Server port")
}

func main() {
	client := NewClient(serverIp, serverPort)
	if client == nil {
		return
	}

	// Read responses from the server
	go client.DealResponse()

	client.Run()
}
