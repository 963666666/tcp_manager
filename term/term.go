package term

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"tcp_manager/codec"
	"tcp_manager/proto"
	"tcp_manager/utils"
	"time"
)

type DevInfo struct {
	AuthKey    string `xorm:"auth_key"`
	Imei       string `xorm:"imei"`
	Vin        string `xorm:"vin"`
	PhoneNum   string `xorm:"pk notnull phone_num"`
	ProvId     uint16 `xorm:"prov_id"`
	CityId     uint16 `xorm:"city_id"`
	Manuf      string `xorm:"manuf"`
	TermType   string `xorm:"term_type"`
	TermId     string `xorm:"term_id"`
	PlateColor int    `xorm:"plate_color"`
	PlateNum   string `xorm:"plate_num"`
}

type GPSData struct {
	Imei      string    `xorm:"pk notnull imei"`
	Stamp     time.Time `xorm:"DataTime pk notnull stamp"`
	WarnFlag  uint32    `xorm:"warn_flag"`
	State     uint32    `xorm:"state"`
	AccState  uint8     `xorm:"acc_state"`
	GpsState  uint8     `xorm:"gps_state"`
	Latitude  uint32    `xorm:"latitude"`
	Longitude uint32    `xorm:"longitude"`
	Altitude  uint16    `xorm:"altitude"`
	Speed     uint16    `xorm:"speed"`
	Direction uint16    `xorm:"direction"'`
}

type Terminal struct {
	authKey   string
	imei      string
	iccid     string
	vin       string
	tboxver   string
	loginTime time.Time
	seqNum    uint16
	phoneNum  []byte
	Conn      net.Conn
	Engine    *xorm.Engine
	Ch        chan int
}

type TermAckBody struct {
	AckSeqNum uint16
	AckId     uint16
	AckResult uint16
}

type PlatAckBody struct {
	AckSeqNum uint16
	AckId     uint16
	AckResult uint16
}

type RegisterBody struct {
	ProId         uint16
	CityId        uint16
	ManufId       []byte `len:"11"`
	TermType      []byte `len:"30"`
	TermId        []byte `len:"30"`
	LicPlateColor uint8
	LicPlate      string
}

type RegisterAckBody struct {
	AckSeqNum uint16
	AckResult uint16
	AuthKey   string
}

type AuthBody struct {
	AuthKeyLen uint8
	AuthKey    string
	Imei       []byte `len:"15"`
	Version    []byte `len:"20"`
}

type GPSInfoBody struct {
	WarnFlag uint32
	State    uint32
	Lat      uint32
	Lng      uint32
	Alt      uint16
	Speed    uint16
	Dir      uint16
	Time     []byte `len:"6"`
}

type CtrlBody struct {
	Cmd   uint8
	Param string
}

func (t *Terminal) Handler(msg proto.Message) []byte {
	if t.phoneNum == nil {
		t.phoneNum = make([]byte, 10)
	}
	copy(t.phoneNum, []byte(msg.HEADER.PhoneNum))
	t.seqNum = msg.HEADER.SeqNum
	switch msg.HEADER.MID {
	case proto.TermAck:
		reqId := utils.Bytes2Word(msg.BODY[2:4])
		if reqId == proto.UpdateReq {

		}
	case proto.Register:
		sql := "insert into banners (`img_url`, `create_time`, `order`, `is_del`) values (?, ?, ?, ?)"
		location := time.Local
		rs, _ := t.Engine.Exec(sql, "i don't know", time.Now().In(location), 1, uint8(1))
		fmt.Println("mysql insert is ", rs)

		devInfo := new(DevInfo)

		devInfo.PhoneNum = strings.TrimLeft(utils.HexBuffToString(t.phoneNum), "0")

		is, _ := t.Engine.Get(devInfo)
		if !is {
			return []byte("con't find this phone number")
		}

		var reg RegisterBody
		_, err := codec.Unmarshal(msg.BODY, &reg)
		if err != nil {
			fmt.Println("err:", err)
		}

		var body []byte
		body, err = codec.Marshal(&RegisterAckBody{
			AckSeqNum: msg.HEADER.SeqNum,
			AckResult: 0,
			AuthKey:   devInfo.AuthKey,
		})
		if err != nil {
			fmt.Println("err: ", err)
		}

		msgAck := &proto.Message{
			HEADER: proto.Header{
				MID:      proto.PlatAck,
				Attr:     proto.MakeAttr(1, false, 0, uint16(len(body))),
				Version:  1,
				PhoneNum: string(t.phoneNum),
				SeqNum:   t.seqNum,
			},
			BODY: body,
		}
		return proto.Packer(msgAck)
	case proto.HeartBeat:
		var err error
		var body []byte
		body, err = codec.Marshal(&PlatAckBody{
			AckSeqNum: msg.HEADER.SeqNum,
			AckId:     msg.HEADER.MID,
			AckResult: 0,
		})
		if err != nil {
			fmt.Println("err: ", err)
		}

		msgAck := &proto.Message{
			HEADER: proto.Header{
				MID:      proto.PlatAck,
				Attr:     proto.MakeAttr(1, false, 0, uint16(len(body))),
				Version:  1,
				PhoneNum: string(t.phoneNum),
				SeqNum:   t.seqNum,
			},
			BODY: body,
		}
		return proto.Packer(msgAck)
	case proto.GpsInfo:
		var gpsInfo GPSInfoBody
		_, err := codec.Unmarshal(msg.BODY, &gpsInfo)
		logrus.Info("gpsInfo is", gpsInfo)
		if err != nil {
			logrus.Println("err: ", err)
		}

		gpsData := &GPSData{
			Imei: t.imei,
			Stamp: time.Now(),
			WarnFlag: gpsInfo.WarnFlag,
			State: gpsInfo.State,
			Latitude: gpsInfo.Lat,
			Longitude: gpsInfo.Lng,
			Altitude: gpsInfo.Alt,
			Speed: gpsInfo.Speed,
			Direction: gpsInfo.Dir,
		}
		if (gpsData.State & 0x00000001) > 0 {
			gpsData.AccState = 1
		} else {
			gpsData.AccState = 0
		}

		if (gpsData.State & 0x00000002) > 0 {
			gpsData.GpsState = 1
		} else  {
			gpsData.GpsState = 0
		}

		_, err = t.Engine.Insert(gpsData)
		if err != nil {
			fmt.Println("insert gps err: ", err)
		}

		var body []byte
		body, err = codec.Marshal(&PlatAckBody{
			AckSeqNum: msg.HEADER.SeqNum,
			AckId: msg.HEADER.MID,
			AckResult:0,
		})
		if err != nil {
			logrus.Println("err: ", err)
		}

		msgAck := &proto.Message{
			HEADER: proto.Header{
				MID: proto.PlatAck,
				Attr: proto.MakeAttr(1, false, 0, uint16(len(body))),
				Version: 1,
				PhoneNum: string(t.phoneNum),
				SeqNum: t.seqNum,
			},
			BODY: body,
		}
		return proto.Packer(msgAck)
	}

	return nil
}
