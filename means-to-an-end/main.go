package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	// Close the connection on function exit
	defer conn.Close()
	raddr := conn.RemoteAddr().String()

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
			log.Fatalln("Could not find bytes to read: Terminating Connection with ", raddr)
			return
		}

		if fault {
			log.Fatalln("Unexpected behaviour from client: Terminating Connection with ", raddr)
			return
		}

		// The 2 variables are 32 bits in size and binary.Read() will try to fill it by reading exactly 4 bytes.
		var x int32
		err = binary.Read(reader, binary.BigEndian, &x)
		// If err == nil, then it mean unexpected behaviour
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				log.Fatalln("Error reading bytes (Unexpected EOF): ", err.Error())
				return
			}
			if err != io.EOF {
				log.Fatalln("Error reading bytes: ", err.Error())
				return
			}
			log.Println("Hit EOF")
			break
		}

		var y int32
		err = binary.Read(reader, binary.BigEndian, &y)
		// If err == nil, then it mean unexpected behaviour
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				log.Fatalln("Error reading bytes (Unexpected EOF): ", err.Error())
				return
			}
			if err != io.EOF {
				log.Fatalln("Error reading bytes: ", err.Error())
				return
			}
			log.Println("Hit EOF")
			break
		}

		// Performing Operations
		if iq {
			log.Println("Performing Insert Operation")
		} else {
			log.Println("Performing Query Operation")
		}
	}

	log.Println("Closing connection with: ", raddr)
}

func main() {
	listener, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalln("Error occured while starting server: ", err.Error())
	}
	defer listener.Close()
	log.Println("Serving on port :2345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error occured while accepting connection: ", err.Error())
		}
		go handleConnection(conn)
	}
}
