package proto

import (
	"fmt"
	"testing"
)

func BenchmarkMakeAttr(b *testing.B) {
	for i := 0; i < b.N; i ++ {
		rs := MakeAttr(byte(1), true, byte(3), uint16(1))
		fmt.Printf("result is %T, value is %v\n", rs, rs)
	}
}