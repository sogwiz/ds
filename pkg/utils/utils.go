package utils

import (
	"context"
	"net"
	"os"
	"os/signal"

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

// SignalCtx returns a context that is cancel when ctrl+c signal is detected
func SignalCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
	return ctx
}
