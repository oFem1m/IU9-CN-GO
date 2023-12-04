package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

func main() {
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Пример простой авторизации. Рекомендуется использовать более безопасные методы.
			if c.User() == "Me" && string(pass) == "123" {
				return nil, nil
			}
			return nil, fmt.Errorf("authentication failed")
		},
	}

	privateBytes, err := ioutil.ReadFile("./lab4/server/key")
	if err != nil {
		fmt.Println("Failed to load private key:", err)
		return
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		fmt.Println("Failed to parse private key:", err)
		return
	}

	config.AddHostKey(private)

	// Запуск SSH сервера
	listener, err := net.Listen("tcp", "localhost:2222")
	if err != nil {
		fmt.Println("Failed to listen on 2222:", err)
		return
	}
	defer listener.Close()

	fmt.Println("SSH server listening on localhost:2222")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept incoming connection:", err)
			continue
		}
		go handleConnection(conn, config)
	}
}

func handleConnection(conn net.Conn, config *ssh.ServerConfig) {
	defer conn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
	if err != nil {
		fmt.Println("Failed to handshake:", err)
		return
	}

	fmt.Println("SSH connection established from", sshConn.RemoteAddr())

	// Обработка запросов
	go ssh.DiscardRequests(reqs)

	// Обработка каналов
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChannel ssh.NewChannel) {
	if t := newChannel.ChannelType(); t != "session" {
		newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
		return
	}

	channel, requests, err := newChannel.Accept()
	if err != nil {
		fmt.Println("Failed to accept channel:", err)
		return
	}
	defer channel.Close()

	for req := range requests {
		switch req.Type {
		case "exec":
			go handleExec(channel, req)
		case "subsystem":
			go handleSubsystem(channel, req)
		default:
			req.Reply(false, nil)
		}
	}
}

func handleExec(channel ssh.Channel, req *ssh.Request) {
	cmd := exec.Command(string(req.Payload[4:]))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command execution failed:", err)
		channel.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		return
	}

	channel.Write(output)
}

func handleSubsystem(channel ssh.Channel, req *ssh.Request) {
	switch string(req.Payload[4:]) {
	case "sftp":
		go handleSFTP(channel)
	default:
		req.Reply(false, nil)
	}
}

func handleSFTP(channel ssh.Channel) {
	serverOptions := []sftp.ServerOption{
		sftp.WithDebug(os.Stderr),
	}

	server, err := sftp.NewServer(
		channel,
		serverOptions...,
	)
	if err != nil {
		fmt.Println("Failed to create SFTP server:", err)
		return
	}
	defer server.Close()

	if err := server.Serve(); err == io.EOF {
		fmt.Println("SFTP session closed by client")
	} else if err != nil {
		fmt.Println("SFTP server exited with error:", err)
	}
}
