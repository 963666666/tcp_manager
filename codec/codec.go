package codec

import (
	"errors"
	"reflect"
	"strconv"
	"tcp_manager/utils"
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
		usedLen = +2
	case reflect.Int32:
		usedLen += 4
	case reflect.Uint32:
		usedLen += 4
	case reflect.Int64:
		usedLen += 8
	case reflect.Uint64:
		usedLen += 8
	case reflect.Float32:
		usedLen += 4
	case reflect.Float64:
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
		for i := 0; i < fieldCount; i++ {
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

	return refMarshal(reflect.ValueOf(v), reflect.StructField{})
}

func refMarshal(v reflect.Value, tag reflect.StructField) ([]byte, error) {
	data := make([]byte, 0)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Int8:
		data = append(data, byte(v.Int()))
	case reflect.Uint8:
		data = append(data, byte(v.Uint()))
	case reflect.Int16:
		temp := utils.Word2Bytes(uint16(v.Int()))
		data = append(data, temp...)
	case reflect.Uint16:
		temp := utils.Word2Bytes(uint16(v.Uint()))
		data = append(data, temp...)
	case reflect.Int32:
		temp := utils.DWord2Bytes(uint32(v.Int()))
		data = append(data, temp...)
	case reflect.Uint32:
		temp := utils.DWord2Bytes(uint32(v.Uint()))
		data = append(data, temp...)
	case reflect.String:
		strLen := tag.Tag.Get("len")
		var lens int = 0
		if strLen == "" {
			lens = v.Len()
		} else {
			lens64, err := strconv.ParseInt(strLen, 10, 0)
			if err != nil {
				return []byte{}, err
			}

			lens = int(lens64)
		}

		if int(lens) > v.Len() {
			zeroSlice := make([]byte, int(lens)-v.Len())
			data = append(data, zeroSlice...)
		}

		data = append(data, v.String()...)
	case reflect.Slice:
		strLen := tag.Tag.Get("len")
		var lens int = 0
		if strLen == "" {
			lens = v.Len()
		} else {
			lens64, err := strconv.ParseInt(strLen, 10, 0)
			if err != nil {
				return []byte{}, err
			}
			lens = int(lens64)
		}

		if int(lens) > v.Len() {
			zeroSlice := make([]byte, int(lens)-v.Len())
			data = append(data, zeroSlice...)
		}
		data = append(data, v.Bytes()...)
	case reflect.Struct:
		fieldCount := v.NumField()

		for i := 0; i < fieldCount; i++ {
			d, err := refMarshal(v.Field(i), v.Type().Field(i))
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
		v.SetInt(int64(utils.Bytes2Word(data)))
		usedLen += 2
	case reflect.Uint16:
		if len(data) < 2 {
			return 0, errors.New("data to short")
		}
		v.SetUint(uint64(utils.Bytes2Word(data)))
		usedLen += 2
	case reflect.Int32:
		if len(data) < 4 {
			return 0, errors.New("data to short")
		}
		v.SetInt(int64(utils.Bytes2DWord(data)))
		usedLen += 4
	case reflect.Uint32:
		if len(data) < 4 {
			return 0, errors.New("data to short")
		}
		v.SetUint(uint64(utils.Bytes2DWord(data)))
		usedLen += 4
	case reflect.Int64:
		v.SetInt(int64(64))
		usedLen += 8
	case reflect.Uint64:
		v.SetUint(uint64(64))
		usedLen += 8
	case reflect.Float32:
		v.SetFloat(32.23)
		usedLen += 4
	case reflect.Float64:
		v.SetFloat(64.46)
		usedLen += 8
	case reflect.String:
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

		for i := 0; i < fieldCount; i++ {
			l, err := refUnmarshal(data[usedLen:], v.Field(i), v.Type().Field(i), streLen)
			if err != nil {
				return 0, err
			}
			usedLen += l
		}

	}
	return usedLen, nil
}
