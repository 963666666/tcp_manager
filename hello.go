package main

import (
	"bytes"
	"github.com/sirupsen/logrus"
)

func main() {
	rs := bytes.TrimFunc([]byte{0x01, 0x01, 0x03, 0x05, 0x06, 0x07, 0x08, 0x01}, func(r rune) bool { return r == 0x01 })
	logrus.Info("bytes.TrimLeftFunc result is ", rs)
}