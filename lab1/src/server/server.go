package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"net"
	"strings"
)

import "file/lab1/src/proto"

// Client - состояние клиента.
type Client struct {
	logger log.Logger        // Объект для печати логов
	conn   *net.TCPConn      // Объект TCP-соединения
	enc    *json.Encoder     // Объект для кодирования и отправки сообщений
	sent   map[string]string //Предложение
}

// NewClient - конструктор клиента, принимает в качестве параметра
// объект TCP-соединения.
func NewClient(conn *net.TCPConn) *Client {
	return &Client{
		logger: log.New(fmt.Sprintf("client %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		sent:   make(map[string]string),
	}
}

// serve - метод, в котором реализован цикл взаимодействия с клиентом.
// Подразумевается, что метод serve будет вызаваться в отдельной go-программе.
func (client *Client) serve() {
	defer client.conn.Close()
	decoder := json.NewDecoder(client.conn)
	for {
		var req proto.Request
		if err := decoder.Decode(&req); err != nil {
			client.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			client.logger.Info("received command", "command", req.Command)
			if client.handleRequest(&req) {
				client.logger.Info("shutting down connection")
				break
			}
		}
	}
}

// handleRequest - метод обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func (client *Client) handleRequest(req *proto.Request) bool {
	switch req.Command {
	case "quit":
		client.respond("ok", nil)
		return true
	case "add":
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			var word proto.Word
			if err := json.Unmarshal(*req.Data, &word); err != nil {
				errorMsg = "malformed data field"
			} else {
				client.logger.Info("adding a word", "value", word.Word, "key", word.Key)
				client.sent[word.Key] = word.Word
			}
		}
		if errorMsg == "" {
			client.respond("ok", nil)
		} else {
			client.logger.Error("adding failed", "reason", errorMsg)
			client.respond("failed", errorMsg)
		}
	case "del":
		if req.Data == nil {
			client.logger.Error("deletion failed", "reason", "key is missing in data field")
			client.respond("failed", "key is missing in data field")
		} else {
			keyToDelete := strings.Trim(string(*req.Data), "\"")
			deletedWord, ok := client.sent[keyToDelete]
			if ok {
				// Удаляем элемент по ключу
				delete(client.sent, keyToDelete)
				client.respond("result", &proto.Word{
					Key:  keyToDelete,
					Word: deletedWord,
				})
			} else {
				client.logger.Error("deletion failed", "reason", "key not found")
				client.respond("failed", "key not found")
			}
		}
	case "sent":
		concatenatedValues := ""
		// Перебираем элементы карты
		for _, value := range client.sent {
			// Конкатенируем значение с предыдущими значениями (если они есть)
			concatenatedValues += value + " "
		}
		client.respond("result", &proto.Word{
			Key:  "sentence",
			Word: concatenatedValues,
		})

	default:
		client.logger.Error("unknown command")
		client.respond("failed", "unknown command")
	}
	return false
}

// respond - вспомогательный метод для передачи ответа с указанным статусом
// и данными. Данные могут быть пустыми (data == nil).
func (client *Client) respond(status string, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&proto.Response{status, &raw})
}

func main() {
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	var addrStr string
	flag.StringVar(&addrStr, "addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	// Разбор адреса, строковое представление которого находится в переменной addrStr.
	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		// Инициация слушания сети на заданном адресе.
		if listener, err := net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
			// Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					// Запуск go-программы для обслуживания клиентов.
					go NewClient(conn).serve()
				}
			}
		}
	}
}
