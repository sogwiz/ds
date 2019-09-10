package main

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

func main() {
	homeDir, _ := os.UserHomeDir()
	fileContent, _ := ioutil.ReadFile(filepath.Join(homeDir, "go", "src", "ds", "cmd", "ds", "main.go"))

	conn, err := net.Dial("tcp", "127.0.0.1:3300")
	if err != nil {
		panic(err)
	}

	_, _ = conn.Write([]byte("user_1/myfile.txt|"))
	_, _ = conn.Write(fileContent)
}
