package utils

import (
	"testing"
)

func TestBytes2DWord(t *testing.T) {
	rs := Bytes2Word([]byte("ab"))
	str := Word2Bytes(rs)
	t.Logf("rs is %v, string is %b", rs, str)
}
