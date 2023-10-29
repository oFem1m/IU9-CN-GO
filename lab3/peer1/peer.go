package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
)

type HashTable struct {
	data  map[string]string
	mutex sync.Mutex
}

type Neighbor struct {
	IP   string
	Port string
}

type Message struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
	Forward bool   `json:"forward"`
}

func (h *HashTable) Add(key, value string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.data[key] = value
}

func (h *HashTable) Delete(key string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.data, key)
}

func (h *HashTable) Find(key string) string {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if val, found := h.data[key]; found {
		return val
	}
	return "Key not found"
}

func main() {
	// Создаем и инициализируем хеш-таблицу
	hashTable := &HashTable{
		data: make(map[string]string),
	}

	// Слайс для хранения информации о соседях
	neighbors := make([]Neighbor, 0)

	// Первый аргумент командной строки - порт, на котором данный пир будет слушать
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run peer.go <port>")
		return
	}

	// Получаем порт из аргументов командной строки
	port := os.Args[1]

	// Запуск слушающего сервера
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Peer listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		go handleConnection(conn, hashTable, &neighbors)
	}
}

func handleConnection(conn net.Conn, hashTable *HashTable, neighbors *[]Neighbor) {
	defer conn.Close()

	// Получение IP и порта соседа
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Connected to peer at %s\n", remoteAddr)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()

		// Распаковка JSON-сообщения
		var msg Message
		err := json.Unmarshal([]byte(command), &msg)
		if err != nil {
			fmt.Printf("Error decoding message: %s\n", err)
			continue
		}

		switch msg.Command {
		case "ADD":
			if msg.Key != "" && msg.Value != "" {
				hashTable.Add(msg.Key, msg.Value)
				fmt.Printf("Added key-value pair: %s-%s\n", msg.Key, msg.Value)
				if !msg.Forward {
					// Отправляем запрос всем соседям, исключая отправителя
					for _, neighbor := range *neighbors {
						if neighbor.IP != msg.IP || neighbor.Port != msg.Port {
							sendRequest(neighbor.IP, neighbor.Port, msg, neighbors)
						}
					}
				}
			}
		case "DELETE":
			if msg.Key != "" {
				hashTable.Delete(msg.Key)
				fmt.Printf("Deleted key: %s\n", msg.Key)
				if !msg.Forward {
					// Отправляем запрос всем соседям, исключая отправителя
					for _, neighbor := range *neighbors {
						if neighbor.IP != msg.IP || neighbor.Port != msg.Port {
							sendRequest(neighbor.IP, neighbor.Port, msg, neighbors)
						}
					}
				}
			}
		case "FIND":
			if msg.Key != "" {
				value := hashTable.Find(msg.Key)
				fmt.Printf("Found value: %s\n", value)
			}
		case "LIST_NEIGHBORS":
			// Вывод информации о соседях в стандартный вывод
			fmt.Println("List of neighbors:")
			for _, neighbor := range *neighbors {
				fmt.Printf("IP: %s, Port: %s\n", neighbor.IP, neighbor.Port)
			}
		case "ADD_NEIGHBOR":
			// Добавляем соседа
			if msg.IP != "" && msg.Port != "" {
				*neighbors = append(*neighbors, Neighbor{IP: msg.IP, Port: msg.Port})
				fmt.Printf("Added neighbor: %s:%s\n", msg.IP, msg.Port)
			}
		}
	}
}

func sendRequest(ip, port string, msg Message, neighbors *[]Neighbor) {
	msg.Forward = true
	msg.IP = ""
	msg.Port = ""
	neighborAddr := ip + ":" + port
	conn, err := net.Dial("tcp", neighborAddr)
	if err != nil {
		fmt.Printf("Error connecting to neighbor %s: %s\n", neighborAddr, err)
		return
	}
	defer conn.Close()
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error encoding message: %s\n", err)
		return
	}
	_, err = conn.Write(append(msgBytes, '\n'))
	if err != nil {
		fmt.Printf("Error sending message to neighbor %s: %s\n", neighborAddr, err)
		return
	}
}
