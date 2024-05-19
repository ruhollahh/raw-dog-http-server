package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var directory string

func main() {
	fmt.Println("starting the server")

	flag.StringVar(&directory, "directory", "./files", "directory of files")
	flag.Parse()

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("failed to bind to port 4221")
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
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	buffer := make([]byte, 1024)
	contentLength, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("error reading response: ", err.Error())
	}

	requestContent := buffer[:contentLength]
	requestHead := strings.Split(string(requestContent), "\r\n\r\n")[0]
	requestBody := requestContent[len(requestHead)+len("\r\n\r\n"):]
	requestHeadParts := strings.Split(requestHead, "\r\n")
	requestLine := requestHeadParts[0]
	requestHeaders := requestHeadParts[1:]

	headers := make(map[string]string)
	for _, header := range requestHeaders {
		keyValuePair := strings.Split(header, ": ")
		headers[strings.ToLower(keyValuePair[0])] = keyValuePair[1]
	}

	requestLineParts := strings.Split(requestLine, " ")
	method := requestLineParts[0]
	route := requestLineParts[1]

	if strings.HasPrefix(route, "/files/") {
		if method == "GET" {
			handleFile(connection, route)
			return
		} else if method == "POST" {
			handleFileUpload(connection, route, requestBody)
			return
		}
	} else if strings.HasPrefix(route, "/echo/") {
		handleEcho(connection, route)
		return
	} else if route == "/user-agent" {
		handleUserAgent(connection, headers["user-agent"])
		return
	} else if route == "/" {
		handleHome(connection)
		return
	}

	handleNotFound(connection)
}
