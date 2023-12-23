package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/smtp"
)

// SMTPSettings содержит настройки для подключения к серверу SMTP
type SMTPSettings struct {
	Server   string
	Port     string
	Username string
	Password string
}

// DBSettings содержит настройки для подключения к базе данных MySQL
type DBSettings struct {
	Host     string
	Database string
	Username string
	Password string
}

// Mail содержит информацию о письме для отправки
type Mail struct {
	To      string
	Subject string
	Body    string
}

// Log содержит информацию для логирования ответов от SMTP-сервера
type Log struct {
	To          string
	Status      string
	Description string
}

func main() {
	// Настройки базы данных
	dbSettings := DBSettings{
		Host:     "students.yss.su",
		Database: "iu9networkslabs",
		Username: "iu9networkslabs",
		Password: "Je2dTYr6",
	}

	// Настройки SMTP
	smtpSettings := SMTPSettings{
		Server:   "mail.nic.ru",
		Port:     "465",
		Username: "dts21@dactyl.su",
		Password: "12345678990DactylSUDTS",
	}

	// Подключение к базе данных
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbSettings.Username, dbSettings.Password, dbSettings.Host, dbSettings.Database))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Подготовка запроса к базе данных
	query := "SELECT username, email, message FROM Volokhov_smtp"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Итерация по результатам запроса
	for rows.Next() {
		var username, email, message string
		if err := rows.Scan(&username, &email, &message); err != nil {
			log.Println(err)
			continue
		}

		// Создание персонализированного HTML-письма
		htmlBody := fmt.Sprintf(`<html><body><p><strong>%s</strong>,</p><p><em>%s</em></p><p>%s</p></body></html>`, username, message, message)

		// Отправка письма
		err := sendMail(smtpSettings, Mail{
			To:      email,
			Subject: "Volokhov Aleksandr IU9-32B",
			Body:    htmlBody,
		}, db)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("Письмо отправлено успешно на %s\n", email)
		fmt.Println("Отправка письма завершена")
	}
}

// sendMail отправляет письмо и сохраняет лог в базе данных
func sendMail(smtpSettings SMTPSettings, mail Mail, db *sql.DB) error {
	// Формирование сообщения
	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html\r\n\r\n%s", mail.To, mail.Subject, mail.Body)

	// Настройка подключения к серверу SMTP
	auth := smtp.PlainAuth("", smtpSettings.Username, smtpSettings.Password, smtpSettings.Server)

	// Подключение к серверу SMTP с использованием TLS
	tlsConfig := &tls.Config{
		ServerName: smtpSettings.Server,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", smtpSettings.Server, smtpSettings.Port), tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpSettings.Server)
	if err != nil {
		return err
	}
	defer client.Close()

	// Аутентификация
	if err := client.Auth(auth); err != nil {
		return err
	}

	// Отправка письма
	if err := client.Mail(smtpSettings.Username); err != nil {
		return err
	}

	if err := client.Rcpt(mail.To); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	// Запись лога в базу данных
	log := Log{
		To:          mail.To,
		Status:      "success",
		Description: "Email sent successfully",
	}
	if err := insertLog(db, log); err != nil {
		fmt.Println("Failed to insert log into the database:", err)
	}

	return nil
}

// insertLog добавляет запись лога в базу данных
func insertLog(db *sql.DB, log Log) error {
	_, err := db.Exec("INSERT INTO Volokhov_logs (to_email, status, description) VALUES (?, ?, ?)", log.To, log.Status, log.Description)
	return err
}
