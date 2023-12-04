package main

import (
	"fmt"
	"net/http"
	"runtime"
	"syscall"

	"github.com/gorilla/websocket"
	"golang.org/x/sys/windows"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type InputData struct {
	PID int `json:"pid"`
}

type OutputData struct {
	ProcessName string `json:"processName,omitempty"`
	Error       string `json:"error,omitempty"`
}

func getProcessName(pid int) (string, error) {
	var processName string

	switch runtime.GOOS {
	case "linux", "darwin":
		// Открываем /proc/<pid>/comm для получения имени процесса в Linux и macOS
		filePath := fmt.Sprintf("/proc/%d/comm", pid)
		file, err := windows.Open(filePath, syscall.O_RDONLY, 0)
		if err != nil {
			return "", err
		}
		defer func(fd windows.Handle) {
			err := windows.Close(fd)
			if err != nil {

			}
		}(file)

		// Читаем имя процесса из файла
		buffer := make([]byte, 1024)
		n, err := windows.Read(file, buffer)
		if err != nil {
			return "", err
		}

		processName = string(buffer[:n])
	case "windows":
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, uint32(pid))
		if err != nil {
			return "", fmt.Errorf("error opening process: %v", err)
		}
		defer func(handle windows.Handle) {
			err := windows.CloseHandle(handle)
			if err != nil {

			}
		}(handle)

		// Выделяем динамическую память для буфера
		bufferSize := windows.MAX_PATH
		buffer := make([]uint16, bufferSize)

		var length uint32

		err = windows.QueryFullProcessImageName(handle, 0, &buffer[0], &length)
		if err != nil {
			return "", fmt.Errorf("error getting process name: %v", err)
		}

		processName = syscall.UTF16ToString(buffer[:length])
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
		var inputData InputData

		// Чтение входных данных от клиента
		err := conn.ReadJSON(&inputData)
		if err != nil {
			fmt.Println(err)
			break
		}

		// Выполнение вычислений
		processName, err := getProcessName(inputData.PID)
		outputData := OutputData{}

		if err != nil {
			// Если произошла ошибка, добавляем ее в OutputData
			outputData.Error = err.Error()
		} else {
			// Если ошибки нет, добавляем имя процесса
			outputData.ProcessName = processName
		}

		// Отправка данных обратно клиенту
		err = conn.WriteJSON(outputData)
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
