package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"log"
)

func main() {
	// Создайте экземпляр парсера RSS
	fp := gofeed.NewParser()

	// Список URL-адресов RSS-каналов, которые вы хотите разобрать
	rssUrls := []string{
		"https://vesti-k.ru/rss/",
		// Другие URL-адреса RSS-каналов, если необходимо
	}

	for _, url := range rssUrls {
		// Разберите RSS-канал
		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Fatalf("Ошибка при разборе RSS с URL %s: %v\n", url, err)
		}

		// Выведите информацию о канале
		fmt.Printf("Title : %s\n", feed.Title)
		fmt.Printf("Description : %s\n", feed.Description)
		fmt.Printf("Number of Items : %d\n", len(feed.Items))

		// Выведите информацию о каждом элементе
		for i, item := range feed.Items {
			fmt.Println()
			fmt.Printf("Item Number : %d\n", i)
			fmt.Printf("Title : %s\n", item.Title)
			fmt.Printf("Link : %s\n", item.Link)
			fmt.Printf("Description : %s\n", item.Description)
			fmt.Printf("Published Date : %s\n", item.Published)
		}
	}
}
