package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	dbUser     = "iu9networkslabs"
	dbPassword = "Je2dTYr6"
	dbName     = "iu9networkslabs"
	dbHost     = "students.yss.su"
)

func main() {
	// Подключение к базе данных MySQL
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Заполнение таблицы данными
	users := []string{"user1", "user2", "user3"}
	emails := []string{"danila@bmstu.posevin.ru", "iu9@bmstu.posevin.ru", "dts21@dactyl.su"}

	for i, user := range users {
		email := emails[i]
		message := fmt.Sprintf("Hello, %s! It's a test message.", user)

		_, err := db.Exec("INSERT INTO Volokhov_smtp (username, email, message) VALUES (?, ?, ?)", user, email, message)
		if err != nil {
			log.Println("Ошибка при вставке данных:", err)
		}
	}

	fmt.Println("Данные успешно добавлены в таблицу Volokhov_smtp.")
}
