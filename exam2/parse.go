package main

import (
	"io"
	"log"
	"net/http"

	"github.com/iotdog/json2table/j2t"
)

func main() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		res, err := http.Get("http://pstgu.yss.su/iu9/mobiledev/lab4_yandex_map/?x=var09")
		if err != nil {
			log.Fatal(err)
		}
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()

		ok, html := j2t.JSON2HtmlTable(string(body), []string{"name", "gps", "address", "tel"}, []string{"title1"})
		if ok {
			log.Print(html)
			w.Write([]byte(html))
		}
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3030", nil))
}
