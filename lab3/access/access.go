package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	var peerAddress string
	var conn net.Conn
	var err error

	// Чтение команд с терминала
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a command: ")
		scanner.Scan()
		command := scanner.Text()

		if command == "exit" {
			// Выход из программы при вводе "exit"
			break
		}

		// Парсинг команды
		parts := strings.Fields(command)
		if len(parts) == 2 && parts[0] == "PEER" {
			peerAddress = parts[1]
			// Устанавливаем соединение с выбранным пиром
			conn, err = net.Dial("tcp", peerAddress)
			if err != nil {
				fmt.Println(peerAddress)
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Println("Connected to peer:", peerAddress)
		} else if conn != nil {
			// Отправляем остальные команды пиру, если соединение установлено
			fmt.Fprintf(conn, command+"\n")
		} else {
			fmt.Println("You must first connect to a peer using 'PEER <port>'")
		}
	}

	if conn != nil {
		conn.Close()
	}
}
