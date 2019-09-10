package utils

import (
	"net"

	"github.com/sirupsen/logrus"
)

// AcceptConn hack to return the conn from a channel.
// We can now wait for a conn in a select statement.
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
