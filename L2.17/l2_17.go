package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	var host string
	var port string
	var timeout int

	flag.StringVar(&host, "h", "", "set host")
	flag.StringVar(&port, "p", "", "set port")
	flag.IntVar(&timeout, "t", 10, "set timeout")
	flag.Parse()

	if host == "" || port == "" {
		log.Println("post и host не должны быть пустыми")
		os.Exit(1)
	}

	conn, err := net.DialTimeout("tcp", host+":"+port, time.Duration(timeout)*time.Second)
	if err != nil {
		log.Printf("Ошибка подключения: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Подключено %s\n", host+":"+port)
	var wg sync.WaitGroup
	// Отправка на сервер
	wg.Go(func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != nil && err != io.EOF {
			log.Printf("ошибка записи: %v\n", err)
		}

		if tcp, ok := conn.(*net.TCPConn); ok {
			log.Printf("Закрываем отправку tcp\n")
			tcp.CloseWrite()
		} else {
			log.Printf("conn не удалось в tcp делаем Close()\n")
			conn.Close()
		}
	})
	// Чтение ответов с сервера
	wg.Go(func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil && err != io.EOF {
			log.Printf("ошибка чтения: %v\n", err)
		}

		conn.Close()
		log.Printf("Закрываем соеденения с %s\n", host+":"+port)
	})

	wg.Wait()
}
