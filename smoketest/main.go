package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	var data []byte
	for {
		buffer := make([]byte, 1024)
		bcount, err := reader.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error: ", err.Error())
				return
			}
			log.Println("Received EOF")
			break
		}
		log.Printf("a connection sent %v bytes of data\n", bcount)
		data = append(data, buffer[:bcount]...)
	}
	log.Println("Received: ", data)

	conn.Write(data)
	fmt.Println("Connection Closed!")
}

func main() {
	listener, err := net.Listen("tcp", ":2345")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on :2345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
