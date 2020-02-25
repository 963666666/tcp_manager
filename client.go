package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"os/signal"
	"tcp_manager/codec"
	"tcp_manager/proto"
	"tcp_manager/term"
)


func main() {
	conn, err := net.Dial("tcp", ":19901")
	if err != nil {
		logrus.Error("net.Dial err:", err)
	}
	defer conn.Close()

	buf := make([]byte, 0)
	for i := 0; i < 3; i ++ {
		buf = append(buf, 0x7e)
		gpsInfoBody := &term.GPSInfoBody{
			WarnFlag: uint32(10),
			State:    uint32(1),
			Lat:      uint32(1),
			Lng:      uint32(3),
			Alt:      uint16(5),
			Speed:    uint16(60),
			Dir:      uint16(7),
			Time:     []byte("20200111"),
		}
		body, err := codec.Marshal(gpsInfoBody)
		if err != nil {
			logrus.Error("codec.Marshal err:", err)
		}

		sendData := &proto.Message{
			HEADER: proto.Header{
				MID: proto.GpsInfo,
				Attr: uint16(len(body)),
				Version: 1,
				PhoneNum: "131000",
				SeqNum: uint16(0),
			},
			BODY: body,
		}
		byteData, err := codec.Marshal(sendData)
		if err != nil {
			logrus.WithFields(logrus.Fields{"hello world": err.Error()}).Error("codec.Marshal err")
		}
		buf = append(buf, byteData...)

		fmt.Println("buf is ", buf)

		_, err = conn.Write(buf)
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	for {
		readBuf := make([]byte, 1024)
		n, err := conn.Read(readBuf)
		if err != nil && err != io.EOF {
			logrus.WithFields(logrus.Fields{"conn.Read err": err.Error()}).Error("err")
			break
		}

		fmt.Printf("hello world len is %d, value is %s\n", n, string(readBuf))


		/*select {
		case sig := <- signalChan:
			logrus.Error("exit sig is ", sig.String())
			break
		default:
			logrus.Info("hello go running")
		}*/
	}
}