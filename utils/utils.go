package utils

func Bytes2Word(data []byte) uint16 {
	if len(data) < 2 {
		return 0
	}
	return (uint16(data[0]) << 8) + uint16(data[1])
}

func Word2Bytes(data uint16) []byte {
	buff := make([]byte, 2)
	buff[0] = byte(data >> 8)
	buff[1] = byte(data)
	return buff
}

func Bytes2DWord(data []byte) uint32 {
	if len(data) < 4 {
		return 0
	}
	return (uint32(data[0]) << 24) + (uint32(data[1]) << 16) + (uint32(data[2]) << 8) + uint32(data[3])
}

func DWord2Bytes(data uint32) []byte {
	buff := make([]byte, 4)
	buff[0] = byte(data >> 24)
	buff[1] = byte(data >> 16)
	buff[2] = byte(data >> 8)
	buff[3] = byte(data)
	return buff
}
