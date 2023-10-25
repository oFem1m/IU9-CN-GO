package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type HashTable struct {
	data  map[string]string
	mutex sync.Mutex
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

	// Слайс для хранения информации о соседях
	neighbors := make([]string, 0)

	// Горутина для чтения команд с командной строки
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			command := scanner.Text()
			if strings.HasPrefix(command, "ADD_PEER") {
				parts := strings.Fields(command)
				if len(parts) != 3 {
					fmt.Println("Invalid ADD_PEER command format")
				} else {
					ip := parts[1]
					peerPort := parts[2]
					neighbors = append(neighbors, ip+":"+peerPort)
					fmt.Printf("Added neighbor: %s\n", ip+":"+peerPort)
				}
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		go handleConnection(conn, hashTable, neighbors)
	}
}

func handleConnection(conn net.Conn, hashTable *HashTable, neighbors []string) {
	defer conn.Close()

	// Получение IP и порта соседа
	remoteAddr := conn.RemoteAddr().String()
	parts := strings.Split(remoteAddr, ":")
	remoteIP := parts[0]
	remotePort := parts[1]

	fmt.Printf("Connected to peer at %s:%s\n", remoteIP, remotePort)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		parts := strings.Fields(command)
		if len(parts) < 2 {
			continue
		}

		switch parts[0] {
		case "ADD":
			if len(parts) == 3 {
				key := parts[1]
				value := parts[2]
				hashTable.Add(key, value)
				fmt.Printf("Added key-value pair: %s-%s\n", key, value)
			}
		case "DELETE":
			if len(parts) == 2 {
				key := parts[1]
				hashTable.Delete(key)
				fmt.Printf("Deleted key: %s\n", key)
			}
		case "FIND":
			if len(parts) == 2 {
				key := parts[1]
				value := hashTable.Find(key)
				fmt.Printf("Found value: %s\n", value)
			}
		case "LIST_NEIGHBORS":
			fmt.Printf("Neighbors: %v\n", neighbors)
		}
	}
}
