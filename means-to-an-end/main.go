package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	raddr := conn.RemoteAddr().String()

	// Close the connection on function exit
	defer conn.Close()

	// Print Remote Address
	log.Println("Received connection from: ", raddr)

	// Create new reader for buffered reading from connection
	reader := bufio.NewReader(conn)

	// Start reading input from the connection
	for {
		// Get Input
		var iq bool    // true is insert, false is query
		fault := false // Set to true if an unexpected behaviour occurs

		op, err := reader.ReadByte()
		switch op {
		case 73:
			iq = true
		case 81:
			iq = false
		default:
			fault = true
		}
		if err != nil {
			log.Println("Could not find bytes to read: Terminating Connection with ", raddr)
			break
		}

		if fault {
			log.Println("Unexpected behaviour from client: Terminating Connection with ", raddr)
			return
		}

		// The 2 variables are 32 bits in size and binary.Read() will try to fill it by reading exactly 4 bytes.
		var x int32
		err = binary.Read(reader, binary.BigEndian, &x)
		if err != nil {
			log.Println("Unexpected behaviour from client: Terminating Connection ", err.Error())
			return
		}

		var y int32
		err = binary.Read(reader, binary.BigEndian, &y)
		if err != nil {
			log.Println("Unexpected behaviour from client: Terminating Connection ", err.Error())
			return
		}

		// Performing Operations
		if iq {
			log.Printf("Performing Insert Operation\tTimestamp: %v\tPrice: %v\n", x, y)
		} else {
			log.Printf("Performing Query Operation\tTimestamp 1: %v\tTimestamp 2: %v\n", x, y)
		}
	}
	log.Println("Closing connection with: ", raddr)
}

func main() {
	listener, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Println("Error occured while starting server: ", err.Error())
	}
	defer listener.Close()
	log.Println("Serving on port :2345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error occured while accepting connection: ", err.Error())
		}
		go handleConnection(conn)
	}
}
