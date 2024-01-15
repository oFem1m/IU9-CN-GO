package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn, ch chan string) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		receivedStr := string(buf[:n])
		ch <- receivedStr
	}
}

func main() {
	listener, err := net.Listen("tcp", "185.139.70.64:0404")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :0303")

	res := ""
	ch := make(chan string)

	go func() {
		for {
			receivedStr := <-ch
			res += receivedStr
			fmt.Println("Result:", res)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, ch)
	}
}
