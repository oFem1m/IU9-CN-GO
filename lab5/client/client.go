package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"strconv"
)

type Process struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"processName,omitempty"`
	Error       string `json:"error,omitempty"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	for {
		// Чтение PID
		fmt.Print("Enter PID: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		pidStr := scanner.Text()

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			fmt.Println("Error: Invalid PID. Please enter a valid number.")
			continue
		}

		// Подготовка данных для отправки серверу
		processData := Process{
			PID: pid,
		}

		// Отправка данных серверу
		err = conn.WriteJSON(processData)
		if err != nil {
			fmt.Println(err)
			break
		}

		// Получение данных от сервера
		var responseData Process
		err = conn.ReadJSON(&responseData)
		if err != nil {
			fmt.Println(err)
			break
		}

		// Проверка наличия ошибки
		if responseData.Error != "" {
			fmt.Printf("Error: %s\n", responseData.Error)
		} else {
			// Вывод результатов
			fmt.Printf("Process Name: %s", responseData.ProcessName)
		}
	}
}
