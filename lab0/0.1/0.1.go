package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
)

// Данные из формы
var data struct {
	Name string
	Age  string
}

// Мьютекс для безопасной работы с данными
var mutex sync.Mutex

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Получаем данные из формы
		name := r.FormValue("name")
		age := r.FormValue("age")

		// Сохраняем данные в глобальной переменной
		mutex.Lock()
		data.Name = name
		data.Age = age
		mutex.Unlock()

		// Перенаправляем пользователя на страницу /data
		http.Redirect(w, r, "/data", http.StatusSeeOther)
		return
	}

	// Отображаем форму для ввода данных
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./lab0/0.1/form.html")
}

func DataHandler(w http.ResponseWriter, r *http.Request) {
	// Защищаем доступ к данным с помощью мьютекса
	mutex.Lock()
	defer mutex.Unlock()

	// Читаем содержимое файла "data.html"
	htmlContent, err := ioutil.ReadFile("./lab0/0.1/data.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Создаем шаблон на основе содержимого файла "data.html"
	tmpl, err := template.New("data").Parse(string(htmlContent))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Отображаем данные из формы в шаблоне и отправляем как ответ клиенту
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", FormHandler)
	http.HandleFunc("/data", DataHandler)

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
