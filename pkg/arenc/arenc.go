// Arithmetic Encoding/Decoding
package arenc

import (
	"fmt"
	"io"
)

// 16 bit encoding
const (
	HIGH     = 0xffff
	HALF     = 0x8000
	THREEQTR = 0xc000
	QTR      = 0x4000
)

type Writer interface {
	io.ByteWriter
	Flush() error
}

type Encoder struct {
	// Writer writes the encoded bits
	Writer io.ByteWriter

	// state for the arithmetic encoder
	lo, hi uint32
	count  uint

	// bits stores up to 32 encoded bits
	bits  uint32
	nBits uint
}

func NewEncoder(w io.ByteWriter) *Encoder {
	return &Encoder{hi: HIGH, Writer: w}
}

func (e *Encoder) EncodeBit(b bool, p0 uint32) error {
	d := e.hi - e.lo + 1
	if b {
		e.lo += (p0 * d) / 0xffff
	} else {
		e.hi = e.lo + (p0*d)/0xffff
	}
	for {
		if e.hi < HALF {
			if err := e.writeBits(false); err != nil {
				return err
			}
		} else if e.lo >= HALF {
			e.lo -= HALF
			e.hi -= HALF
			if err := e.writeBits(true); err != nil {
				return err
			}
		} else if QTR <= e.lo && e.hi < THREEQTR {
			e.count++
			e.lo -= QTR
			e.hi -= QTR
		} else {
			break
		}
		e.lo *= 2
		e.hi = 2*e.hi + 1
	}
	return nil
}

// Writes output of the encoder to underlying Writer in little endian order
func (e *Encoder) writeBits(b bool) error {
	for e.count > 0 {
		e.count--
		if b {
			e.bits = 2*e.bits + 1
		} else {
			e.bits = 2 * e.bits
		}
		e.nBits++
		for e.nBits >= 8 {
			fmt.Printf("%016b %016b %08b\n", e.lo, e.hi, byte(e.bits))
			if err := e.Writer.WriteByte(byte(e.bits)); err != nil {
				return err
			}
			e.bits >>= 8
			e.nBits -= 8
		}
	}
	return nil
}

type Decoder struct {
	// state for the arithmetic decoder
	lo, hi, code uint32

	// bits stores the decoded bits
	bits  uint32
	nBits uint

	// encBits stores the encoded bits that the decoder reads from
	encBits  uint32
	nEncBits uint
}

func (dc *Decoder) DecodeBit(p0 uint32) bool {
	d := dc.hi - dc.lo + 1
	r := (dc.code - dc.lo + 1) << 16 / d
	b := r >= p0

	for {
		if dc.hi < HALF {
		} else if dc.lo >= HALF {
			dc.lo -= HALF
			dc.hi -= HALF
			dc.code -= HALF
		} else if QTR <= dc.lo && dc.hi < THREEQTR {
			dc.code -= QTR
			dc.lo -= QTR
			dc.hi -= QTR
		} else {
			break
		}
		dc.lo *= 2
		dc.hi = 2*dc.hi + 1
		dc.code = 2*dc.code + dc.encBits&1
		dc.encBits >>= 1
	}
	return b
}

//func (e *Encoder) Read(buf []byte) (n int, err error){
//}
