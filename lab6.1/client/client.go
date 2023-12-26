package main

import (
	"bufio"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"path"
	"strings"
)

const (
	ftpHost  = "students.yss.su"
	ftpLogin = "ftpiu8"
	ftpPass  = "3Ru7yOTA"
)

func main() {
	// Подключение к FTP-серверу
	client, err := ftp.Dial(ftpHost + ":21")
	if err != nil {
		panic(err)
	}
	defer client.Quit()

	err = client.Login(ftpLogin, ftpPass)
	if err != nil {
		panic(err)
	}

	fmt.Println("FTP client connected successfully!")

	// Запуск интерактивного режима
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Введите команду (upload, download, create, delete, list, exit): ")
		scanner.Scan()
		command := scanner.Text()

		switch strings.ToLower(command) {
		case "upload":
			fmt.Print("Введите локальный путь к файлу: ")
			scanner.Scan()
			localPath := scanner.Text()
			fmt.Print("Введите удаленный путь для сохранения файла: ")
			scanner.Scan()
			remotePath := scanner.Text()
			uploadFile(client, localPath, remotePath)

		case "download":
			fmt.Print("Введите удаленный путь к файлу: ")
			scanner.Scan()
			remotePath := scanner.Text()
			fmt.Print("Введите локальный путь для сохранения файла: ")
			scanner.Scan()
			localPath := scanner.Text()
			downloadFile(client, remotePath, localPath)

		case "create":
			fmt.Print("Введите имя новой директории: ")
			scanner.Scan()
			directoryName := scanner.Text()
			createDirectory(client, directoryName)

		case "delete":
			fmt.Print("Введите удаленный путь к файлу для удаления: ")
			scanner.Scan()
			remotePath := scanner.Text()
			deleteFile(client, remotePath)

		case "list":
			fmt.Print("Введите удаленный путь: ")
			scanner.Scan()
			remotePath := scanner.Text()
			listDirectory(client, remotePath)

		case "exit":
			fmt.Println("Выход из программы.")
			return

		default:
			fmt.Println("Неверная команда. Пожалуйста, введите корректную команду.")
		}
	}
}

// Функция загрузки файла на FTP-сервер
func uploadFile(client *ftp.ServerConn, localPath, remotePath string) {
	file, err := os.Open(localPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = client.Stor(remotePath+"/"+getFileName(localPath), file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s uploaded successfully.\n", localPath)
}

// Функция скачивания файла с FTP-сервера
func downloadFile(client *ftp.ServerConn, remotePath, localPath string) {
	resp, err := client.Retr(remotePath)
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	file, err := os.Create(localPath + "/" + getFileName(remotePath))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("File %s downloaded successfully.\n", remotePath)
}

// Функция создания директории на FTP-сервере
func createDirectory(client *ftp.ServerConn, directoryName string) {
	err := client.MakeDir(directoryName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Directory %s created successfully.\n", directoryName)
}

// Функция удаления файла на FTP-сервере
func deleteFile(client *ftp.ServerConn, remotePath string) {
	err := client.Delete(remotePath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("File %s deleted successfully.\n", remotePath)
}

// Функция получения содержимого директории на FTP-сервере
func listDirectory(client *ftp.ServerConn, directoryPath string) {
	entries, err := client.List(directoryPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Directory listing for %s:\n", directoryPath)
	for _, entry := range entries {
		fmt.Println(entry.Name)
	}
}

// Функция для получения имени файла из полного пути
func getFileName(filePath string) string {
	_, file := path.Split(filePath)
	return file
}
