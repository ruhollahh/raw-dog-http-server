package main

import (
	"fmt"
	"net"
	"os"
)

func handleNotFound(connection net.Conn) {
	_, err := connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	if err != nil {
		fmt.Println("error writing response: ", err.Error())
		os.Exit(1)
	}

	err = connection.Close()
	if err != nil {
		fmt.Println("error closing connection: ", err.Error())
		os.Exit(1)
	}
}

func handleHome(connection net.Conn) {
	_, err := connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	if err != nil {
		fmt.Println("error writing response: ", err.Error())
		os.Exit(1)
	}

	err = connection.Close()
	if err != nil {
		fmt.Println("error closing connection: ", err.Error())
		os.Exit(1)
	}
}

func handleEcho(connection net.Conn, route string) {
	message := route[len("/echo/"):]
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(message), message)
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("error writing response: ", err.Error())
		os.Exit(1)
	}

	err = connection.Close()
	if err != nil {
		fmt.Println("error closing connection: ", err.Error())
		os.Exit(1)
	}
}

func handleUserAgent(connection net.Conn, userAgent string) {
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("error writing response: ", err.Error())
		os.Exit(1)
	}

	err = connection.Close()
	if err != nil {
		fmt.Println("error closing connection: ", err.Error())
		os.Exit(1)
	}
}
