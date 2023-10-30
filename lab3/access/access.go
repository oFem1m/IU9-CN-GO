package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Forward bool   `json:"forward"`
}

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
		if parts[0] == "PEER" && len(parts) == 2 {
			peerAddress = "localhost:" + parts[1]
			// Устанавливаем соединение с выбранным пиром
			conn, err = net.Dial("tcp", peerAddress)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}
			fmt.Println("Connected to peer:", peerAddress)
		} else if parts[0] == "ADD_NEIGHBOR" && len(parts) == 3 {
			// Подготовка сообщения в формате JSON
			msg := Message{
				Command: "ADD_NEIGHBOR",
				IP:      parts[1],
				Port:    parts[2],
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error encoding message: %s\n", err)
				continue
			}

			// Отправка JSON-сообщения
			_, err = conn.Write(append(msgBytes, '\n'))
			if err != nil {
				fmt.Printf("Error sending message: %s\n", err)
			}
		} else if conn != nil {
			// Отправляем остальные команды пиру, если соединение установлено

			// Подготовка сообщения в формате JSON
			msg := Message{
				Command: parts[0],
				Key:     "",
				Value:   "",
			}

			if len(parts) > 1 {
				msg.Key = parts[1]
			}

			if len(parts) > 2 {
				msg.Value = parts[2]
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("Error encoding message: %s\n", err)
				continue
			}

			// Отправка JSON-сообщения
			_, err = conn.Write(append(msgBytes, '\n'))
			if err != nil {
				fmt.Printf("Error sending message: %s\n", err)
			}
		} else {
			fmt.Println("You must first connect to a peer using 'PEER <port>'")
		}
	}

	if conn != nil {
		err := conn.Close()
		if err != nil {
			return
		}
	}
}
