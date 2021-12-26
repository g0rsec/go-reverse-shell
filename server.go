package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var infoLogger *log.Logger = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.LUTC|log.Lshortfile)
var errorLogger *log.Logger = log.New(os.Stdout, "[ERROR] ", log.LstdFlags|log.LUTC|log.Lshortfile)

func handleErr(err error, errorLogger *log.Logger) {
	if err != nil {
		errorLogger.Fatal(strings.Title(err.Error()))
	}
}

func getData(conn net.Conn) string {
	d, err := bufio.NewReader(conn).ReadBytes('\n')
	handleErr(err, errorLogger)
	decodedData, _ := base64.StdEncoding.DecodeString(strings.Split(strings.TrimSpace(string(d)), "\n")[0])
	return string(decodedData)
}

func sendData(conn net.Conn, data string) {
	encodedData := base64.StdEncoding.EncodeToString([]byte(data)) + "\n"
	_, err := conn.Write([]byte(encodedData))
	handleErr(err, errorLogger)
	// infoLogger.Printf("Sent %v bytes of data to remote client %v.", n, conn.RemoteAddr().String())
}

func main() {
	infoLogger.Print("Starting server...")
	ln, err := net.Listen("tcp", "127.0.0.1:8000")
	handleErr(err, errorLogger)
	defer ln.Close()
	infoLogger.Printf("Listening on %v/%v...", ln.Addr().String(), ln.Addr().Network())

	for {
		conn, err := ln.Accept()
		handleErr(err, errorLogger)
		defer conn.Close()
		infoLogger.Printf("Client %v connected.", conn.RemoteAddr().String())
		decodedData := getData(conn)
		if decodedData == "goRS-ID" {
			infoLogger.Print("Client ID given verified.")
			go func() {
				for {
					fmt.Print(getData(conn))
				}
			}()
			for {
				reader := bufio.NewReader(os.Stdin)
				cmd, _ := reader.ReadString('\n')
				sendData(conn, cmd)
			}
		} else {
			infoLogger.Print("Wrong client ID given. Closing connection.")
			sendData(conn, "Wrong client ID given. Closing connection.")
			conn.Close()
		}
	}
}
