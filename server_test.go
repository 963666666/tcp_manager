package main

import (
	"fmt"
	"net"
	"tcp_manager/codec"
	"tcp_manager/proto"
	"testing"
	"time"
)

func BenchmarkTcpServer(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		go func() {
			conn, err := net.Dial("tcp", "192.168.1.109:19903")
			if err != nil {
				return
			}

			defer conn.Close()

			n, err := conn.Write([]byte("123456"))
			if err != nil {
				return
			}
			fmt.Println("n is ", n)
			for {
				time.Sleep(10 * time.Second)
			}
		}()
	}
	b.Log("b.N is ", b.N)
}

func TestTcpServer(t *testing.T) {
	conn, err := net.Dial("tcp", "192.168.1.109:19901")
	if err != nil {
		t.Error("net.Dial err:", err)
	}
	defer conn.Close()


		buf := make([]byte, 1024)

		t.Logf("read len is %v\n", buf)


		sendData := &proto.Message{
			HEADER: proto.Header{
				MID: proto.Register,
				Attr: uint16(0),
				Version: 1,
				PhoneNum: "131000",
			},
			BODY: []byte{},
		}
		byteData, err := codec.Marshal(sendData)

		_, err = conn.Write(byteData)
		t.Logf("hello world first: %v", byteData)


		//fmt.Printf("byteData is %v\n", byteData)
		//var uSendData = &proto.Message{}
		//lens, err := codec.Unmarshal(byteData, uSendData)
		//fmt.Printf("hello lens is %d, data is %v\n", lens, uSendData)
		//
		//	n, err = conn.Write(byteData)
		//	t.Logf("hello world first: %v", n)

}
