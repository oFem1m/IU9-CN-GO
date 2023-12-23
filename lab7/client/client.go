package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	// Ввод данных с клавиатуры
	fmt.Print("Введите адрес получателя (To): ")
	to := readInput()

	fmt.Print("Введите тему сообщения (Subject): ")
	subject := readInput()

	fmt.Print("Введите текст сообщения (Message body): ")
	messageBody := readInput()

	// Аутентификационные данные
	username := "dts21@dactyl.su"
	password := "12345678990DactylSUDTS"

	// SMTP-сервер и порт
	smtpServer := "mail.nic.ru"
	smtpPort := 465

	// Использование SSL
	useSSL := true

	// Формирование темы сообщения с указанием фамилии, имени и группы студента
	messageSubject := fmt.Sprintf("%s", subject)

	// Формирование тела сообщения
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, messageSubject, messageBody)

	// Подготовка конфигурации TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !useSSL,
		ServerName:         smtpServer,
	}

	// Установка соединения с сервером
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpServer, smtpPort), tlsConfig)
	if err != nil {
		fmt.Println("Ошибка при установке соединения:", err)
		return
	}
	defer conn.Close()

	// Аутентификация на сервере
	auth := smtp.PlainAuth("", username, password, smtpServer)

	// Установка клиентского сеанса
	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		fmt.Println("Ошибка при создании клиентского сеанса:", err)
		return
	}
	defer client.Quit()

	// Аутентификация
	if err := client.Auth(auth); err != nil {
		fmt.Println("Ошибка при аутентификации:", err)
		return
	}

	// Отправка письма
	if err := client.Mail(username); err != nil {
		fmt.Println("Ошибка при указании отправителя:", err)
		return
	}

	if err := client.Rcpt(to); err != nil {
		fmt.Println("Ошибка при указании получателя:", err)
		return
	}

	writer, err := client.Data()
	if err != nil {
		fmt.Println("Ошибка при отправке данных:", err)
		return
	}
	defer writer.Close()

	_, err = writer.Write([]byte(message))
	if err != nil {
		fmt.Println("Ошибка при записи данных:", err)
		return
	}

	fmt.Println("Письмо успешно отправлено.")
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
