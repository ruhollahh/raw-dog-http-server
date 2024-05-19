package main

import (
	"errors"
	"fmt"
	"net"
	"os"
)

func handleFileUpload(connection net.Conn, route string, fileContent []byte) {
	fileName := route[len("/files/"):]
	filePath := fmt.Sprintf("%s/%s", directory, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("error creating file: ", err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error closing file: ", err.Error())
		}
	}(file)

	_, err = file.Write(fileContent)
	if err != nil {
		fmt.Println("error writing to file: ", err.Error())
		return
	}

	response := "HTTP/1.1 201 Created\r\n\r\n"
	_, err = connection.Write([]byte(response))
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

func handleFile(connection net.Conn, route string) {
	fileName := route[len("/files/"):]
	filePath := fmt.Sprintf("%s/%s", directory, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			handleNotFound(connection)
			return
		}
		fmt.Println("error reading file: ", err.Error())
		return
	}

	responseText := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n", len(content))
	response := append([]byte(responseText), content...)
	_, err = connection.Write(response)
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
