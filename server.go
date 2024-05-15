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
		headersLine := requestLines[1 : len(requestLines)-1-1]
		headers := make(map[string]string)
		for _, header := range headersLine {
			keyValuePair := strings.Split(header, ": ")
			headers[strings.ToLower(keyValuePair[0])] = keyValuePair[1]
		}
		route := strings.Split(requestLines[0], " ")[1]

		if strings.HasPrefix(route, "/echo/") {
			handleEcho(con, route)
		} else if route == "/user-agent" {
			handleUserAgent(con, headers["user-agent"])
		} else if route == "/" {
			handleHome(con)
		} else {
			handleNotFound(con)
		}
	}
}
