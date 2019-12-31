package codec

import (
	"debug/dwarf"
	"errors"
	"go/types"
	"reflect"
	"strconv"
)

func RequireLen(v interface{}) (int, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, errors.New("error")
	}

	return refRequireLen(rv, reflect.StructField{})
}

func refRequireLen(value reflect.Value, tag reflect.StructField) (int, error) {
	var usedLen int = 0
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Int8:
		usedLen += 1
	case reflect.Uint8:
		usedLen += 1
	case reflect.Int16:
		usedLen += 2
	case reflect.Uint16:
		usedLen = + 2
	case reflect.Int32:
		usedLen += 4
	case reflect.Uint32:
		usedLen += 4
	case reflect.Int64:
		usedLen += 8
	case reflect.Uint64:
		usedLen += 8
	case reflect.String:
		strLen := tag.Tag.Get("len")
		if strLen == "" {
			return 0, nil
		}
		lens, err := strconv.ParseInt(strLen, 10, 0)
		if err != nil {
			return 0, err
		}
		usedLen += int(lens)
	case reflect.Slice:
		strLen := tag.Tag.Get("len")
		if strLen == "" {
			return 0, nil
		}
		lens, err := strconv.ParseInt(strLen, 10, 0)
		if err != nil {
			return 0, err
		}
		usedLen += int(lens)
	case reflect.Struct:
		fieldCount := value.NumField()
		for i := 0; i < fieldCount; i ++ {
			l, err := refRequireLen(value.Field(i), value.Type().Field(i))
			if err != nil {
				return 0, nil
			}

			usedLen += l
		}
	}
	return usedLen, nil
}

func Marshal(v interface{}) ([]byte, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return []byte{}, errors.New("error")
	}

	return refMarshal(rv, reflect.StructField{})
}

func refMarshal(value reflect.Value, field reflect.StructField) ([]byte, error) {
	data := make([]byte, 0)
	if value.Kind() == reflect.Ptr {
		value.Elem()
	}

	switch value.Kind() {
	case reflect.Int8:
		data = append(data, byte(value.Int()))
	case reflect.Uint8:
		data = append(data, byte(value.Uint()))
	case reflect.Int16:
		temp := Word2Bytes(uint16(value.Int()))
		data = append(data, temp...)
	case reflect.Uint16:
		temp := Word2Bytes(uint16(value.Uint()))
		data = append(data, temp...)
	case reflect.Int32:
		temp := DWord2Bytes(uint32(value.Int()))
		data = append(data, temp...)
	case reflect.Uint32:
		temp := DWord2Bytes(uint32(value.Int()))
		data = append(data, temp...)
	case reflect.String:
		strLen := field.Tag.Get("len")
		lens, err := strconv.ParseInt(strLen, 10, 0)
		if err != nil {
			return []byte{}, err
		}
		if int(lens) > value.Len() {
			zeroSlize := make([]byte, int(lens)- value.Len())
			data = append(data, zeroSlize...)
		}
		data = append(data, value.String()...)
	case reflect.Slice:
		strLen := field.Tag.Get("len")
		lens, err := strconv.ParseInt(strLen, 10, 0)
		if err != nil {
			return []byte{}, nil
		}

		if int(lens) > value.Len() {
			zeroSlize := make([]byte, int(lens) - value.Len())
			data = append(data, zeroSlize...)
		}
		data = append(data, value.Bytes()...)
	case reflect.Struct:
		fieldCount := value.NumField()
		for i := 0; i < fieldCount; i ++ {
			d, err := refMarshal(value.Field(i), value.Type().Field(i))
			if err != nil {
				return []byte{}, err
			}
			data = append(data, d...)
		}
	}
	return data, nil
}

func Unmarshal(data []byte, v interface{}) (int, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return 0, errors.New("type err")
	}
	lens, err := RequireLen(v)
	if err != nil {
		return 0, err
	}
	if len(data) < lens {
		return 0, errors.New("data too short")
	}

	return refUnmarshal(data, rv, reflect.StructField{}, len(data)-lens)
}

func refUnmarshal(data []byte, v reflect.Value, tag reflect.StructField, streLen int) (int, error) {
	var usedLen int = 0
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Int8:
		v.SetInt(int64(data[0]))
		usedLen += 1
	case reflect.Uint8:
		v.SetUint(uint64(data[0]))
		usedLen += 1
	case reflect.Int16:
		if len(data) < 2 {
			return 0, errors.New("data to short")
		}
		v.SetInt(int64(Bytes2Word(data)))
		usedLen += 2
	case reflect.Uint16:
		if len(data) < 2 {
			return 0, errors.New("data to short")
		}
		v.SetUint(uint64(Bytes2Word(data)))

		usedLen += 2
	case reflect.Int32:
		if len(data) < 4 {
			return 0, errors.New("data to short")
		}
		usedLen += 4
	case reflect.Uint32:
		if len(data) < 4 {
			return 0, errors.New("data to short")
		}
		v.SetUint(uint64(Bytes2Word(data)))
	case reflect.Int64:
		v.SetInt(64)
		usedLen += 8
	case reflect.Uint64:
		v.SetUint(64)
		usedLen += 8
	case reflect.Float32:
		v.SetFloat(32.23)
		usedLen += 4
	case reflect.Float64:
		v.SetFloat(64.46)
		usedLen += 8
	case reflect.String:
		strLen := tag.Tag.Get("len")
		var lens = 0
		if strLen == "" {
			lens = streLen
		} else {
			lens64, err := strconv.ParseInt(strLen, 10, 0)
			if err != nil {
				return 0, err
			}
			lens = int(lens64)
		}

		if len(data) < int(lens) {
			return 0, errors.New("data to short")
		}
		v.SetString(string(data[:lens]))

		usedLen += lens

	case reflect.Slice:
		strLen := tag.Tag.Get("len")
		var lens int = 0
		if strLen == "" {
			lens = streLen
		} else {
			lens64, err := strconv.ParseInt(strLen, 10, 0)
			if err != nil {
				return 0, err
			}
			lens = int(lens64)
		}
		v.SetBytes(data[:lens])
		usedLen += int(lens)

	case reflect.Struct:
		fieldCount := v.NumField()

		for i := 0; i < fieldCount; i ++ {
			l, err := refUnmarshal(data[usedLen:], v.Field(i), v.Type().Field(i), streLen)
			if err != nil {
				return 0, err
			}
			usedLen += l
		}

	}
	return usedLen, nil
}

func Bytes2Word(data []byte) uint16 {
	if len(data) < 2 {
		return 0
	}
	return (uint16(data[0])<<8 + uint16(data[1]))
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
	return (uint32(data[0]<<24) + (uint32(data[1]) << 16) + uint32(data[2])<<8) + uint32(data[3]))
}

func DWord2Bytes(data uint32) []byte {
	buff := make([]byte, 4)
	buff[0] = byte(data >> 24)
	buff[1] = byte(data >> 16)
	buff[2] = byte(data >> 8)
	buff[3] = byte(data)
	return buff
}
