package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	if username == "" {
		fmt.Println("Username cannot be empty.")
		os.Exit(1)
	}

	// Send username to server
	_, err = conn.Write([]byte(username))
	if err != nil {
		fmt.Println("Error sending username:", err)
		os.Exit(1)
	}

	fmt.Println("Welcome to the chatroom, " + username + "!")
	fmt.Println("Type 'exit' to quit or press Ctrl+C.")
	fmt.Println("------------------------------------")

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\nExiting chatroom...")
		conn.Close()
		os.Exit(0)
	}()

	// Goroutine to receive messages from server
	go func() {
		reader := bufio.NewReader(conn)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				// Connection closed
				return
			}
			fmt.Print(strings.TrimSpace(message) + "\n")
		}
	}()

	// Main loop - send messages to server
	for {
		scanner.Scan()
		message := strings.TrimSpace(scanner.Text())

		if message == "exit" {
			fmt.Println("Exiting chatroom...")
			break
		}

		if message == "" {
			continue
		}

		// Print own message locally (no self-echo from server)
		fmt.Printf("You: %s\n", message)

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			fmt.Println("Server may be down. Exiting...")
			break
		}
	}
}