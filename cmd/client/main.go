package main

import (
	"io"
	"net"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func main() {
	masterHostname := "127.0.0.1:3300"

	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, "go", "src", "ds", "cmd", "ds", "main.go")
	f, err := os.Open(path)
	if err != nil {
		logrus.Fatal("file not found", err)
	}
	defer f.Close()

	conn, err := net.Dial("tcp", masterHostname)
	if err != nil {
		panic(err)
	}

	_, _ = conn.Write([]byte("user_1/myfile.txt|"))
	_, _ = io.Copy(conn, f)
}
