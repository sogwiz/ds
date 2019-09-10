package utils

import (
	"net"

	"github.com/sirupsen/logrus"
)

func AcceptConn(l net.Listener) <-chan net.Conn {
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
