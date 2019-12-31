package codec

import "testing"

func TestUnmarshal(t *testing.T) {
	type Body struct {
		Age1 int8
		Age2 int16
	}

	data := []byte{0x01, 0x02, 0x03}
	pack := &Body{127, 32767}
	i, err := Unmarshal(data, pack)
	if err != nil {
		t.Errorf("err: %s\n", err.Error())
	}

	t.Logf("i is: %v", i)
	t.Log("pack: ", t)
}
