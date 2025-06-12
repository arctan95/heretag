package main

func main() {
	server := NewServer("0.0.0.0", 6969)
	server.Start()
}
