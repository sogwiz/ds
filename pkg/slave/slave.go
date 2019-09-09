package slave

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func handleRequest(conn net.Conn) {
	// TODO: read... all the data (maybe more than 1024 bytes ?)
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		logrus.Error("Error reading:", err.Error())
	}

	fmt.Println("received:", string(buf))

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