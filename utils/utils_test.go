package utils

import (
	"testing"
)

func TestHexToAsc(t *testing.T) {
	rs := HexBuffToString([]byte{0xaa, 0xcc})
	t.Logf("rs is %v", rs)
}
