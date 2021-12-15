package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net"
	"os"
	"unicode/utf8"
)

var logger *log.Logger

func mylog(msg string) {
	fmt.Println(msg)
	logger.Println(msg)
}

func myerr(msg string) {
	fmt.Println(msg)
	logger.Fatalf(msg)
}

func localAddr() net.IP {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func encode(s string) string {
	r := ""
	for _, l := range s {
		c, _ := utf8.DecodeRuneInString(string(l + 3))
		r += string(c)
	}
	return r
}

func decode(s string) string {
	r := ""
	for _, l := range s {
		c, _ := utf8.DecodeRuneInString(string(l - 3))
		r += string(c)
	}
	return r
}

func send(address, msg string) {
	conn, err := net.Dial("udp", address+":12345")
	if err != nil {
		myerr("Не удалось отправить сообщение! " + err.Error())
		return
	}
	fmt.Fprintf(conn, msg)
}

type clients struct {
	id        int    `json:"id"`
	sfm       string `json:"sfm"`
	residence string `json:"residance"`
}

type librarians struct {
	id  int    `json:"id"`
	sfm string `json:"sfm"`
}

type films struct {
	name     string `json:"name"`
	year     string `json:"year"`
	director string `json:"director"`
	genre    string `json:"genre"`
	timeline int    `json:"timeline"`
	studio   string `json:"studio"`
}

type cassettes struct {
	id    int     `json:"id"`
	price float32 `json:"price"`
	film  string  `json:"film"`
	year  int     `json:"year"`
}

type givings struct {
	id       int `json:"id"`
	client   int `json:"client"`
	cassette int `json:"cassette"`
	issued   int `json:"issued"`
}

func main() {
	f, err := os.Open("log.txt")
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов!", err)
	}
	defer f.Close()
	logger = log.New(f, "", log.LstdFlags)
	defer mylog("Сервер отключен.")
	mylog("Сервер запущен.")

	pc, err := net.ListenPacket("udp", localAddr().String()+":12345")
	if err != nil {
		myerr(err.Error())
	}
	defer pc.Close()
	for {
		var r string
		buf := make([]byte, 1024)
		_, address, err := pc.ReadFrom(buf)
		if err != nil {
			myerr("Ошибка подключения " + address.String() + ", " + err.Error())
			continue
		}
		msg := decode(string(buf))
		mylog("Подключился " + address.String())
		mylog("Запрос " + msg)
		db, err := sql.Open("sqlite3", "db.db")
		if err != nil {
			myerr("Ошибка получения данных из базы " + err.Error())
		}
		querys, err := db.Query(msg)
		if err != nil {
			myerr(err.Error())
		}
		var res []*[]*sql.NullString
		i := 0
		for querys.Next() {
			var t []*sql.NullString
			res = append(res, &t)
			querys.ScanString(res[i])
			i++
		}
		for _, e := range res {
			for _, ee := range *e {
				r += (*ee).String + " "
			}
			r += "\n"
		}
		send(address.String(), r)
		err = querys.Close()
		if err != nil {
			myerr(err.Error())
		}
		db.Close()
	}
}
