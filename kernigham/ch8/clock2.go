package main

import (
	"fmt"
	"time"
	"log"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	fmt.Printf("%#v\n", listener)
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		fmt.Printf("%#v, %#v, %v\n", count, conn, conn)
		count += 1
		go handleConn(conn) // handle one connection at a time
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}
