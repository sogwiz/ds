package main

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, "go", "src", "ds", "cmd", "ds", "main.go")
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Fatal("file not found", err)
	}

	conn, err := net.Dial("tcp", "127.0.0.1:3300")
	if err != nil {
		panic(err)
	}

	_, _ = conn.Write([]byte("user_1/myfile.txt|"))
	_, _ = conn.Write(fileContent)
}
