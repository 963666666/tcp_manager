package proto

import (
	"bytes"
	"errors"
	"fmt"
	"tcp_manager/utils"
)

const (
	ProtoHeader byte = 0x7e

	TermAck     uint16 = 0x0001
	Register    uint16 = 0x0100
	RegisterAck uint16 = 0x8100
	Unregister  uint16 = 0x0003
	Login       uint16 = 0x0102
	HeartBeat   uint16 = 0x0002
	GpsInfo     uint16 = 0x0200
	PlatAck     uint16 = 0x8001
	UpdateReq   uint16 = 0x8108
	CtrlReq     uint16 = 0x8105
)

type Header struct {
	MID       uint16
	Attr      uint16
	Version   uint8
	PhoneNum  string
	SeqNum    uint16
	MultiFlag MultiFlag
}

type MultiFlag struct {
	MsgSum   uint16
	MsgIndex uint16
}

type Message struct {
	HEADER Header
	BODY   []byte
}


func (h *Header) IsMulti() bool {
	if ((h.Attr >> 12) & 0x0001) > 0 {
		return true
	}
	return false
}

func (h *Header) BodyLen() int {
	return int(h.Attr & 0x03ff)
}

func MakeAttr(verFlag byte, mut bool, enc byte, lens uint16) uint16 {
	attr := lens & 0x03FF

	if verFlag > 0 {
		attr = attr & 0x4000
	}

	if mut {
		attr = attr & 0x2000
	}

	encMask := (uint16(enc) & 0x0007) << 10
	return attr + encMask
}

func Version() string {
	return "1.0.0"
}

func Name() string {
	return "jtt808"
}

func Filter(data []byte) ([]Message, int, error) {
	var usedLen int = 0
	msgList := make([]Message, 0)
	var cnt int = 0
	for {
		cnt ++
		if cnt > 10 {
			return []Message{}, 0, fmt.Errorf("time too much")
		}
		if usedLen > len(data) {
			break
		}
		msg, lens, err := filterSigle(data[usedLen:])
		if err != nil {
			usedLen += lens
			fmt.Println("err: ", err.Error())
			return msgList, usedLen, nil
		}
		usedLen += lens
		msgList = append(msgList, msg)
	}
	return msgList, usedLen, nil
}

func filterSigle(data []byte) (Message, int, error) {
	var usedLen int = 0
	startIndex := bytes.IndexByte(data, ProtoHeader)
	if startIndex >= 0 {
		usedLen = startIndex + 1
		endIndex := bytes.IndexByte(data[usedLen:], ProtoHeader)
		if endIndex >= 0 {
			msg, err := frameParser(data[startIndex+1 : endIndex])
			if err != nil {
				return Message{}, endIndex, err
			}
			return msg, endIndex + 1, nil

		}
		return Message{}, startIndex, errors.New("can't find end flag")
	}
	return Message{}, len(data), errors.New("can't find start flag")
}

func Escape(data, oldBytes, newBytes []byte) []byte {
	buff := make([]byte, 0)

	var startIndex int = 0
	for startIndex < len(data) {
		index := bytes.Index(data[startIndex:], oldBytes)

		if index > 0 {
			buff = append(buff, data[startIndex:index]...)
			buff = append(buff, newBytes...)
			startIndex = index + len(oldBytes)
		} else {
			buff = append(buff, data[startIndex:]...)
			startIndex = len(data)
		}
	}
	return buff
}

func frameParser(data []byte) (Message, error) {
	if (len(data) + 2) < (17 + 3) {
		return Message{}, errors.New("header is too short")
	}

	// 不包含帧头尾
	frameData := Escape(data[:len(data)], []byte{0x7d, 0x02}, []byte{0x7e})
	frameData = Escape(frameData, []byte{0x7d, 0x01}, []byte{0x7d})

	rawcs := CheckSum(frameData[:len(frameData)-1])
	if rawcs != frameData[len(frameData)-1] {
		return Message{}, fmt.Errorf("cs is not match: %d--%d", rawcs, frameData[len(frameData)-1])
	}

	var usedLen int = 0
	var msg Message
	msg.HEADER.MID = utils.Bytes2Word(frameData[usedLen:])
	usedLen += 2
	msg.HEADER.Attr = utils.Bytes2Word(frameData[usedLen:])
	usedLen += 2
	msg.HEADER.Version = frameData[usedLen]
	usedLen += 1

	tempPhone := bytes.TrimLeftFunc(frameData[usedLen:usedLen+10], func(r rune) bool { return r == 0x00 })
	msg.HEADER.PhoneNum = string(tempPhone)
	usedLen += 10
	msg.HEADER.SeqNum = utils.Bytes2Word(frameData[usedLen:])
	usedLen += 2
	if msg.HEADER.IsMulti() {
		msg.HEADER.MultiFlag.MsgSum = utils.Bytes2Word(frameData[usedLen:])
		usedLen += 2
		msg.HEADER.MultiFlag.MsgIndex = utils.Bytes2Word(frameData[usedLen:])
		usedLen += 2
	}

	if len(frameData) < usedLen {
		return Message{}, fmt.Errorf("flag code is too short")
	}
	msg.BODY = make([]byte, len(frameData)-usedLen)
	copy(msg.BODY, frameData[usedLen:len(frameData)])
	usedLen = len(frameData)

	return msg, nil
}

func CheckSum(data []byte) byte {
	var sum byte = 0
	for _, itemData := range data {
		sum ^= itemData
	}
	return sum
}

func Packer(msg *Message) []byte {
	data := make([]byte, 0)
	tempBytes := utils.Word2Bytes(msg.HEADER.MID)
	data = append(data, tempBytes...)
	dataLen := uint16(len(msg.BODY)) & 0x03FF
	dataLen = dataLen | 0x4000

	tempBytes = utils.Word2Bytes(dataLen)
	data = append(data, tempBytes...)

	data = append(data, msg.HEADER.Version)

	if len(msg.HEADER.PhoneNum) < 10 {
		data = append(data, make([]byte, 10-len(msg.HEADER.PhoneNum))...)
		data = append(data, msg.HEADER.PhoneNum...)
	} else {
		data = append(data, msg.HEADER.PhoneNum[:10]...)
	}

	tempBytes = utils.Word2Bytes(msg.HEADER.SeqNum)
	data = append(data, tempBytes...)

	if msg.HEADER.IsMulti() {
		data = append(data, utils.Word2Bytes(msg.HEADER.MultiFlag.MsgSum)...)
		data = append(data, utils.Word2Bytes(msg.HEADER.MultiFlag.MsgIndex)...)
	}

	data = append(data, msg.BODY...)

	// 添加头尾
	var tmpData []byte = []byte{0x7e}

	for _, item := range data {
		if item == 0x7d {
			tmpData = append(tmpData, 0x7d, 0x01)
		} else if item == 0x7e {
			tmpData = append(tmpData, 0x7d, 0x02)
		} else {
			tmpData = append(tmpData, item)
		}
	}
	tmpData = append(tmpData, 0x7e)
	return tmpData
}
