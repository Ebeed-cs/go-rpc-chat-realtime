package main

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	conn     net.Conn
	username string
	outgoing chan string
}

type ChatServer struct {
	mu      sync.RWMutex
	clients map[net.Conn]*Client
}

func NewChatServer() *ChatServer {
	return &ChatServer{
		clients: make(map[net.Conn]*Client),
	}
}

func (s *ChatServer) addClient(conn net.Conn, username string) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	client := &Client{
		conn:     conn,
		username: username,
		outgoing: make(chan string, 10),
	}
	s.clients[conn] = client

	return client
}

func (s *ChatServer) removeClient(conn net.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, ok := s.clients[conn]; ok {
		close(client.outgoing)
		delete(s.clients, conn)
	}
}

func (s *ChatServer) broadcast(message string, sender net.Conn) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for conn, client := range s.clients {
		if conn != sender {
			select {
			case client.outgoing <- message:
			default:
				// channel full, skip this client
			}
		}
	}
}

func (s *ChatServer) handleClient(conn net.Conn) {
	defer conn.Close()

	// Read username first
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading username:", err)
		return
	}
	username := string(buffer[:n])

	client := s.addClient(conn, username)
	defer s.removeClient(conn)

	// Notify all other clients about new user
	joinMsg := fmt.Sprintf("User %s joined", username)
	s.broadcast(joinMsg, conn)

	// Start sender goroutine
	go func() {
		for msg := range client.outgoing {
			_, err := conn.Write([]byte(msg + "\n"))
			if err != nil {
				return
			}
		}
	}()

	// Receiver loop - read messages from this client
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			// Client disconnected
			leaveMsg := fmt.Sprintf("User %s left", username)
			s.broadcast(leaveMsg, conn)
			return
		}

		message := string(buffer[:n])
		fullMsg := fmt.Sprintf("%s: %s", username, message)
		s.broadcast(fullMsg, conn)
	}
}

func main() {
	server := NewChatServer()

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listener.Close()

	fmt.Println("Chat server started on :1234")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go server.handleClient(conn)
	}
}