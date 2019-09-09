package slave

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func handleRequest(conn net.Conn) {
	reader := bufio.NewReader(conn)
	filename, _ := reader.ReadString('|')
	filename = strings.TrimSuffix(filename, "|")

	// Read the incoming connection into the buffer.
	dir, file := filepath.Split(filename)
	_ = os.MkdirAll(filepath.Join("data", dir), 0777)
	fo, err := os.Create(filepath.Join("data", dir, file))
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		}
	}

	fmt.Println("received a file: " + filename)

	if err := conn.Close(); err != nil {
		logrus.Error(err)
	}
}

func acceptConn(l net.Listener) <-chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		conn, err := l.Accept()
		if err != nil {
			logrus.Error(err)
			close(ch)
		}
		ch <- conn
		close(ch)
	}()
	return ch
}

func StartTCPServer() {
	c1, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})

	go func(ctx context.Context) {
		logrus.Info("slave server listening on localhost:3333")
		l, err := net.Listen("tcp", "localhost:3333")
		if err != nil {
			panic(err)
		}
		for {
			select {
			case <-ctx.Done():
				logrus.Info("cancelled")
				close(exitCh)
				return
			case conn := <-acceptConn(l):
				go handleRequest(conn)
			}
		}
	}(c1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	<-exitCh
}
