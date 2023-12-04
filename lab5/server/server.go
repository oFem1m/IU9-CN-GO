package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Process struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"processName,omitempty"`
	Error       string `json:"error,omitempty"`
}

func getProcessName(pid int) (string, error) {
	var processName string

	switch runtime.GOOS {
	case "linux", "darwin":
		commPath := filepath.Join("/proc", strconv.Itoa(pid), "comm")
		data, err := ioutil.ReadFile(commPath)
		if err != nil {
			return "", err
		}
		processName = string(data)
	default:
		return "", fmt.Errorf("unsupported operating system")
	}

	return processName, nil
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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
		var process Process

		// Чтение входных данных от клиента
		err := conn.ReadJSON(&process)
		if err != nil {
			fmt.Println(err)
			break
		}

		processName, err := getProcessName(process.PID)
		if err != nil {
			// Если произошла ошибка, добавляем ее в OutputData
			process.Error = err.Error()
		} else {
			// Если ошибки нет, добавляем имя процесса
			process.ProcessName = processName
		}

		// Отправка данных обратно клиенту
		err = conn.WriteJSON(process)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/", handleConnection)
	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
