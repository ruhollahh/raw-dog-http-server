package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("starting the server")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("error closing listener: ", err.Error())
			os.Exit(1)
		}
	}(listener)

	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err.Error())
			os.Exit(1)
		}

		buffer := make([]byte, 1024)
		contentLength, err := con.Read(buffer)
		if err != nil {
			fmt.Println("error reading response: ", err.Error())
			os.Exit(1)
		}

		requestContent := string(buffer[:contentLength])
		requestLines := strings.Split(requestContent, "\r\n")
		route := strings.Split(requestLines[0], " ")[1]

		var response string

		if strings.HasPrefix(route, "/echo/") {
			message := route[len("/echo/"):]
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)
		} else if route == "/" {
			response = "HTTP/1.1 200 OK\r\n\r\n"
		} else {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		}

		_, err = con.Write([]byte(response))
		if err != nil {
			fmt.Println("error writing response: ", err.Error())
			os.Exit(1)
		}

		err = con.Close()
		if err != nil {
			fmt.Println("error closing connection: ", err.Error())
			os.Exit(1)
		}
	}
}
