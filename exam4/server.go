package main

import (
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	message := string(buffer[:n])
	fmt.Printf("Received message from client: %s\n", message)

	// Преобразование в заглавные буквы
	response := strings.ToUpper(message)

	// Отправка обратно клиенту
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	fmt.Printf("Sent response to client: %s\n", response)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
