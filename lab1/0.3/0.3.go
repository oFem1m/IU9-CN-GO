package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"html/template"
	"net/http"
	"sync"
)

// Данные из формы
var data struct {
	RssUrl string
}

// Мьютекс для безопасной работы с данными
var mutex sync.Mutex

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Получаем URL RSS-канала из формы
		rssUrl := r.FormValue("rssUrl")

		// Сохраняем URL в глобальной переменной
		mutex.Lock()
		data.RssUrl = rssUrl
		mutex.Unlock()

		// Перенаправляем пользователя на страницу /news
		http.Redirect(w, r, "/news", http.StatusSeeOther)
		return
	}

	// Отображаем форму для ввода URL RSS-канала
	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, "./0.3/form.html")
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	// Защищаем доступ к данным с помощью мьютекса
	mutex.Lock()
	rssUrl := data.RssUrl
	mutex.Unlock()

	// Парсим RSS-канал
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssUrl)
	if err != nil {
		http.Error(w, "Ошибка при разборе RSS-канала", http.StatusInternalServerError)
		return
	}

	// Отображаем ссылки на новости с использованием шаблона news.html
	tmpl, err := template.ParseFiles("./0.3/news.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, feed)
	if err != nil {
		http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", FormHandler)
	http.HandleFunc("/news", NewsHandler) // Добавляем новый путь "/news" для вывода новостей

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
