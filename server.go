package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"tcp_manager/proto"
	"tcp_manager/term"
)

func CheckError(err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{"Error:": err.Error()}).Error("check")
		logrus.WithFields(logrus.Fields{"Error:": err.Error()}).Error("check")
		os.Exit(1)
	}
}

func recvConnMsg(conn net.Conn) {

	buf := make([]byte, 0)
	addr := conn.RemoteAddr()
	logrus.WithFields(logrus.Fields{"network": addr.Network(), "ip": addr.String()}).Info("recv")

	var t *term.Terminal = &term.Terminal{
		Conn: conn,
		Ch: make(chan int),
		Engine: Engine,
	}

	ipAddress := addr.String()
	defer conn.Close()

	for {
		tempBuf := make([]byte, 1024)
		n, err := conn.Read(tempBuf)

		if err != nil {
			logrus.WithFields(logrus.Fields{"network": addr.Network(), "ip": addr.String()}).Info("closed")
			return
		}

		buf = append(buf, tempBuf[:n]...)
		var outLog string
		for _, val := range buf {
			outLog += fmt.Sprintf("%02X", val)
		}

		logrus.WithFields(logrus.Fields{"data": outLog}).Info("<--- ")
		msg, lens, err := proto.Filter(buf)
		if err != nil {
			logrus.WithFields(logrus.Fields{"network": addr.Network(), "ip": ipAddress}).Info("proto.Filter err")
		}

		logrus.WithFields(logrus.Fields{"msg is ": msg}).Info("unMarshal msg is ")

		buf = buf[:lens]

		for len(msg) > 0 {
			sendBuf := t.Handler(msg[0])
			conn.Write([]byte("recv success"))

			logrus.WithFields(logrus.Fields{"sendBuf is ": sendBuf}).Info("sendBuf is")
			if sendBuf != nil {
				outLog = ""
				for _, val := range sendBuf {
					outLog += fmt.Sprintf("%02X", val)
				}
				logrus.WithFields(logrus.Fields{"data": outLog}).Info("---> ")
				conn.Write(sendBuf)

				msg = msg[1:]
			}
		}
	}
}

var Engine *xorm.Engine


func main() {
	var err error
	Engine, err = xorm.NewEngine("mysql", "root:123456@/test?charset=utf8")
	CheckError(err)

	defer Engine.Close()

	listener, err := net.Listen("tcp", ":19901")
	CheckError(err)
	defer listener.Close()

	for {
		newConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go recvConnMsg(newConn)
	}
}
