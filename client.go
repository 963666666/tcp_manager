package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"tcp_manager/codec"
	"tcp_manager/proto"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.1.109:19901")
	if err != nil {
		logrus.Error("net.Dial err:", err)
	}
	defer conn.Close()


		buf := make([]byte, 0)
		buf = append(buf, 0x7e)

	for {
		sendData := &proto.Message{
			HEADER: proto.Header{
				MID: proto.Register,
				Attr: uint16(3),
				Version: 1,
				PhoneNum: "131000",
				SeqNum: uint16(0),
			},
			BODY: []byte{0x01, 0x01, 0x01},
		}
		byteData, err := codec.Marshal(sendData)
		if err != nil {
			logrus.WithFields(logrus.Fields{"hello world": err.Error()}).Error("codec.Marshal err")
		}
		buf = append(buf, byteData...)
		buf = append(buf, 0x7e)

		fmt.Println("buf is ", buf)

		_, err = conn.Write(buf)
		fmt.Printf("hello world first: %v", len(buf))
	}




	//fmt.Printf("byteData is %v\n", byteData)
	//var uSendData = &proto.Message{}
	//lens, err := codec.Unmarshal(byteData, uSendData)
	//fmt.Printf("hello lens is %d, data is %v\n", lens, uSendData)
	//
	//	n, err = conn.Write(byteData)
	//	t.Logf("hello world first: %v", n)
}