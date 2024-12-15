package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"math"
	"net"
)

type PrimeRequest struct {
	Method *string                `json:"method"`
	Number *float64               `json:"number"`
	Extras map[string]interface{} `json:"-"`
}

type PrimeResponse struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func isPrime(n float64) bool {
	if n <= 1 {
		return false
	}
	if n != float64(int(n)) {
		return false
	}
	if n <= 3 {
		return true
	}
	if math.Mod(n, 2) == 0 || math.Mod(n, 3) == 0 {
		return false
	}
	for i := 5; i*i <= int(n); i += 6 {
		if math.Mod(n, float64(i)) == 0 || math.Mod(n, float64(i)+2) == 0 {
			return false
		}
	}
	return true
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	log.Println("New connection from: ", conn.RemoteAddr().String())
	for {
		line, err := reader.ReadBytes(10)

		var req PrimeRequest
		var res PrimeResponse

		if err != nil {
			if err != io.EOF {
				log.Println("An Error Occurred while processing connection: ", err.Error())
				return
			}
			break
		}
		log.Println("Got Request: ", string(line))
		err = json.Unmarshal(line, &req)

		// Performing Validations
		if err != nil {
			log.Printf("Error occured with request %v\n", err.Error())
			_, err := writer.Write([]byte("Malformed Request\n"))
			if err != nil {
				log.Fatalln("Error: ", err.Error())
			}
			writer.Flush()
			log.Println("Connection Closed due to Malformed Request")
			return
		}

		if req.Number == nil || req.Method == nil {
			log.Println("Malformed Request, missing required fields!")
			_, err := writer.Write([]byte("Malformed Request\n"))
			if err != nil {
				log.Fatalln("Error: ", err.Error())
			}
			writer.Flush()
			log.Println("Connection Closed due to Malformed Request")
			return
		}

		if *req.Method != "isPrime" {
			log.Println("Malformed Request method is not 'isPrime'")
			_, err := writer.Write([]byte("Malformed Request\n"))
			if err != nil {
				log.Fatalln("Error: ", err.Error())
			}
			writer.Flush()
			log.Println("Connection Closed due to Malformed Request")
			return
		}

		res.Method = "isPrime"
		res.Prime = isPrime(*req.Number)

		response, err := json.Marshal(res)
		if err != nil {
			log.Println("Error occurred during response marshalling: ", err.Error())
		}

		response = append(response, 10)
		log.Println(string(response))

		_, err = writer.Write(response)
		if err != nil {
			log.Fatalln("Error: ", err.Error())
		}
		err = writer.Flush()
		if err != nil {
			log.Fatalln(err.Error())
		}

	}
	log.Println("Connection Closed")
}

func main() {

	listener, err := net.Listen("tcp", ":2345")
	if err != nil {
		log.Fatalln("An Error Occured: ", err.Error())
	}
	defer listener.Close()

	log.Println("Server is now listening on port :2345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("An Error Occurred: ", err.Error())
		}
		go handleConnection(conn)
	}

}
