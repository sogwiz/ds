package slave

import (
	"fmt"
	"net"
)

func handleRequest(conn net.Conn) {
	// TODO: read... all the data (maybe more than 1024 bytes ?)
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	fmt.Println("received:", string(buf))

	if err := conn.Close(); err != nil {
		fmt.Println(err)
	}
}

func StartTCPServer() {
	//c1, cancel := context.WithCancel(context.Background())
	//exitCh := make(chan struct{})
	//
	//go func(ctx context.Context) {
	l, err := net.Listen("tcp", "localhost:3333")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleRequest(conn)
	}
	//}(c1)
	//
	//signalCh := make(chan os.Signal, 1)
	//signal.Notify(signalCh, os.Interrupt)
	//go func() {
	//	select {
	//	case <-signalCh:
	//		cancel()
	//		return
	//	}
	//}()
	//<-exitCh
}
