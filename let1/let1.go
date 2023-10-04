package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	passwordURL := "http://pstgu.yss.su/iu9/networks/let1/getkey.php"

	name := "Волохов Александр"

	md5Hash := md5.Sum([]byte(name))
	hashStr := hex.EncodeToString(md5Hash[:])

	fullURL := passwordURL + "?hash=" + hashStr

	response, err := http.Get(fullURL)
	if err != nil {
		fmt.Println("Ошибка при запросе пароля:", err)
		return
	}
	defer response.Body.Close()

	passwordBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}
	password := strings.TrimSpace(string(passwordBytes))

	password = password[6:]

	sendURL := "http://pstgu.yss.su/iu9/networks/let1/send_from_go.php?subject=let1_ИУ9-32Б_Волохов_Александр&fio=Волохов_Александр&pass=" + password

	response, err = http.Get(sendURL)
	if err != nil {
		fmt.Println("Ошибка при отправке данных:", err)
		return
	}
	defer response.Body.Close()

	fmt.Println("Данные успешно отправлены с паролем:", password)
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}
	fmt.Println("Ответ сервера:\n", string(responseBytes))
}
