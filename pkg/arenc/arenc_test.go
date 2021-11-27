package arenc

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestEncodeBits(t *testing.T) {
	f, err := os.OpenFile("testing.txt", os.O_WRONLY, 666)
	if err != nil {
		fmt.Println(err)
	}
	bw := bufio.NewWriterSize(f, 1)
	e := NewEncoder(bw)
	var p0 uint32 = HALF
	for i := 0; i < 32; i++ {
		//01110110
		e.EncodeBit(false, p0)
		e.EncodeBit(true, p0)
		e.EncodeBit(true, p0)
		e.EncodeBit(true, p0)
		e.EncodeBit(false, p0)
		e.EncodeBit(true, p0)
		e.EncodeBit(true, p0)
		e.EncodeBit(false, p0)
	}
}
