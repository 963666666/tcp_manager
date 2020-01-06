package main

import (
	"fmt"
	"net"
)

type Terminal struct {
	authKey string
	imei string
	iccid string
	vin string
	tboxver string
	phoneNum string
	Conn net.Conn
}

var connManager map[string]*Terminal

func recvConnMsg(conn net.Conn) {
	addr := conn.RemoteAddr()

	var term = &Terminal{
		Conn: conn,
	}

	connManager[addr.String()] = term

	for {
		tempbuf := make([]byte, 1024)
		n, err := conn.Read(tempbuf)
		if err != nil {
			return
		}
		fmt.Println("rcv: ", tempbuf[:n])
	}
}

func TcpServer(addr string) {
	connManager = make(map[string]*Terminal)
	listenSock, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	fmt.Println("hello tcp server is running")

	for {
		newConn, err := listenSock.Accept()
		if err != nil {
			continue
		}

		go recvConnMsg(newConn)
	}
}

func main() {
	//TcpServer(":19903")
	fmt.Println(string(0x7e))
}