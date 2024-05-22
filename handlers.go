package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
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

func handleEcho(connection net.Conn, route string, reqHeaders map[string]string) {
	message := route[len("/echo/"):]

	acceptEncoding := reqHeaders["accept-encoding"]
	var contentEncodingHeader string
	var resBodyBuffer bytes.Buffer
	if strings.Contains(acceptEncoding, "gzip") {
		contentEncodingHeader = fmt.Sprintf("Content-Encoding: %s\r\n", "gzip")

		gz := gzip.NewWriter(&resBodyBuffer)
		_, err := gz.Write([]byte(message))
		if err != nil {
			fmt.Println("error encoding response to gzip", err)
		}
		err = gz.Close()
		if err != nil {
			fmt.Println("error closing gzip writer", err)
		}
	} else {
		resBodyBuffer.Write([]byte(message))
	}

	resBody := resBodyBuffer.Bytes()
	resHeaders := fmt.Sprintf("Content-Type: text/plain\r\nContent-Length: %d\r\n%s", len(resBody), contentEncodingHeader)

	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n%s\r\n", resHeaders)
	_, err := connection.Write(append([]byte(response), resBody...))
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
