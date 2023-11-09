package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("SSH Server Address: ")
	serverAddress, _ := reader.ReadString('\n')
	serverAddress = strings.TrimSpace(serverAddress)

	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", serverAddress, config)
	if err != nil {
		fmt.Println("Failed to connect to the SSH server:", err)
		return
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {

		}
	}(client)

	currentDir := "" // Текущая директория

	for {
		fmt.Print("Enter SSH command (or 'exit' to quit): ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		if command == "exit" {
			break
		}
		session, err := client.NewSession()
		if err != nil {
			fmt.Println("Failed to create SSH session:", err)
			continue
		}

		switch {
		case strings.HasPrefix(command, "ls"):
			output, err := session.CombinedOutput("ls " + currentDir)
			err = session.Close()
			if err != nil {
				return
			}
			if err != nil {
				fmt.Println("Failed to list directory:", err)
			} else {
				fmt.Println("Directory contents:\n", string(output))
			}

		case strings.HasPrefix(command, "cd"):
			newDir := strings.TrimSpace(strings.TrimPrefix(command, "cd"))
			Path := newDir
			if currentDir != "" {
				Path = currentDir + "/" + newDir
			}
			output, err := session.CombinedOutput("cd " + Path + " && pwd")
			currentDir = Path
			err = session.Close()
			if err != nil {
				return
			}
			if err != nil {
				fmt.Println("Failed to change directory:", err)
			} else {
				currentDir = strings.TrimSpace(string(output))
				fmt.Println("Current directory:", currentDir)
			}

		case strings.HasPrefix(command, "mkdir"):
			dirName := strings.TrimSpace(strings.TrimPrefix(command, "mkdir"))
			_, err := session.CombinedOutput("mkdir " + currentDir + "/" + dirName)
			if err != nil {
				fmt.Println("Failed to create directory:", err)
			} else {
				fmt.Println("Directory created:", dirName)
			}
		case strings.HasPrefix(command, "rmdir"):
			dirName := strings.TrimSpace(strings.TrimPrefix(command, "rmdir"))
			Path := dirName
			if currentDir != "" {
				Path = currentDir + "/" + dirName
			}
			_, err := session.CombinedOutput("rmdir " + Path)
			if err != nil {
				fmt.Println("Failed to remove directory:", err)
			} else {
				fmt.Println("Directory removed:", dirName)
			}
		case strings.HasPrefix(command, "mv"):
			args := strings.Split(command, " ")
			if len(args) != 3 {
				fmt.Println("Usage: mv source destination")
			} else {
				sourcePath := args[1]
				destPath := args[2]
				_, err := session.CombinedOutput("mv " + sourcePath + " " + destPath)
				if err != nil {
					fmt.Println("Failed to move:", err)
				} else {
					fmt.Println("Moved:", sourcePath, "to", destPath)
				}
			}
		default:
			output, err := session.CombinedOutput(command)

			if err != nil {
				fmt.Println("Failed to execute the command:", err)
			} else {
				fmt.Println("Command output:\n", string(output))
			}
		}
	}

	fmt.Println("SSH session ended.")
}
