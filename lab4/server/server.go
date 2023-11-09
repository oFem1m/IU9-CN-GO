package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gliderlabs/ssh"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		command := strings.Fields(s.Command()[0])

		switch command[0] {
		case "listDir":
			listDir(s, command[1:])
		case "createDir":
			createDir(s, command[1:])
		case "removeDir":
			removeDir(s, command[1:])
		case "moveFile":
			moveFile(s, command[1:])
		case "removeFile":
			removeFile(s, command[1:])
		case "runCommand":
			runCommand(s, command[1:])
		default:
			fmt.Fprintln(s, "Unknown command")
		}
	})

	err := ssh.ListenAndServe("localhost:2222", nil)
	if err != nil {
		fmt.Println("Failed to start SSH server:", err)
	}
}

func createDir(s ssh.Session, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(s, "Usage: createDir <dirName>")
		return
	}

	dirName := args[0]
	err := os.Mkdir(dirName, 0755)
	if err != nil {
		fmt.Fprintln(s, "Failed to create directory:", err)
		return
	}
	fmt.Fprintln(s, "Directory created successfully")
}

func removeDir(s ssh.Session, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(s, "Usage: removeDir <dirName>")
		return
	}

	dirName := args[0]
	err := os.RemoveAll(dirName)
	if err != nil {
		fmt.Fprintln(s, "Failed to remove directory:", err)
		return
	}
	fmt.Fprintln(s, "Directory removed successfully")
}

func listDir(s ssh.Session, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(s, "Usage: listDir <dirName>")
		return
	}

	dirName := args[0]
	files, err := os.ReadDir(dirName)
	if err != nil {
		fmt.Fprintln(s, "Failed to list directory:", err)
		return
	}

	for _, file := range files {
		fmt.Fprintln(s, file.Name())
	}
}

func moveFile(s ssh.Session, args []string) {
	if len(args) != 2 {
		fmt.Fprintln(s, "Usage: moveFile <srcPath> <destPath>")
		return
	}

	srcPath := args[0]
	destPath := args[1]
	err := os.Rename(srcPath, destPath)
	if err != nil {
		fmt.Fprintln(s, "Failed to move file:", err)
		return
	}
	fmt.Fprintln(s, "File moved successfully")
}

func removeFile(s ssh.Session, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(s, "Usage: removeFile <fileName>")
		return
	}

	fileName := args[0]
	err := os.Remove(fileName)
	if err != nil {
		fmt.Fprintln(s, "Failed to remove file:", err)
		return
	}
	fmt.Fprintln(s, "File removed successfully")
}

func runCommand(s ssh.Session, args []string) {
	cmd := strings.Join(args, " ")

	result, err := executeCommand(cmd)
	if err != nil {
		fmt.Fprintln(s, "Error executing command:", err)
		return
	}

	io.WriteString(s, result)
}

func executeCommand(command string) (string, error) {
	return "", fmt.Errorf("execution of external commands is not implemented in this example")
}
