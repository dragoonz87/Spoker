package frame

import (
	"bufio"
	"fmt"
	"log"
	"math"
)

const (
	FINBIT     = 0x1 << 7
	RSV1BIT    = 0x1 << 6
	RSV2BIT    = 0x1 << 5
	RSV3BIT    = 0x1 << 4
	OPCODEBITS = 0xf
	MASKBIT    = 0x1 << 7
	PLLENBITS  = 0x7f
)

func read(r *bufio.Reader, n int) ([]byte, error) {
	bytes, err := r.Peek(n)
	if err != nil {
		log.Printf("couldn't peek %d\n", n)
		return nil, err
	}

	r.Discard(n)

	return bytes, nil
}

type FrameError struct {
	arg int
}

func (fe *FrameError) Error() string {
	return fmt.Sprintf("%d", fe.arg)
}

type Frame struct {
	isFinal                 bool
	rsv1                    bool
	rsv2                    bool
	rsv3                    bool
	Opcode                  uint8
	isMasked                bool
	payloadLengthExtensions uint8
	payloadLength           uint64
	maskingKey              []byte
	PayloadData             []byte
}

func New(isFinal bool, data []byte) *Frame {
	payloadLength := len(data)
	var extension uint8
	if payloadLength > 125 {
		if payloadLength <= math.MaxUint16 {
			extension = 1
		} else {
			extension = 2
		}
	} else {
		extension = 0
	}

	return &Frame{
		isFinal:                 isFinal,
		rsv1:                    false,
		rsv2:                    false,
		rsv3:                    false,
		Opcode:                  1,
		isMasked:                false,
		payloadLengthExtensions: extension,
		payloadLength:           uint64(payloadLength),
		maskingKey:              nil,
		PayloadData:             data,
	}
}

func NewString(isFinal bool, data string) *Frame {
	dataBytes := []byte{}
	dataBytes = append(dataBytes, data...)
	return New(isFinal, dataBytes)
}

func ReadToNewFrame(r *bufio.Reader) (*Frame, error) {
	bytes, err := read(r, 2)
	if err != nil {
		return nil, err
	}

	fin := bytes[0]&FINBIT != 0
	rsv1 := bytes[0]&RSV1BIT != 0
	rsv2 := bytes[0]&RSV2BIT != 0
	rsv3 := bytes[0]&RSV3BIT != 0
	opcode := bytes[0] & OPCODEBITS
	masked := bytes[1]&MASKBIT != 0
	payloadlen := bytes[1] & PLLENBITS

	if !fin {
		log.Println("not final")
		return nil, &FrameError{arg: 1}
	}

	if rsv1 {
		log.Println("rv1 enabled")
		return nil, &FrameError{arg: 2}
	}

	if rsv2 {
		log.Println("rv2 enabled")
		return nil, &FrameError{arg: 3}
	}

	if rsv3 {
		log.Println("rv3 enabled")
		return nil, &FrameError{arg: 4}
	}

	if !masked {
		log.Println("not masked")
		return nil, &FrameError{arg: 5}
	}

	var extension uint64
	var extended uint8

	if payloadlen == 126 {
		bytes, err := read(r, 2)
		if err != nil {
			return nil, err
		}

		extension = uint64(bytes[0])<<8 + uint64(bytes[1])
		extended = 1
	} else if payloadlen == 127 {
		bytes, err := read(r, 8)
		if err != nil {
			return nil, err
		}

		extension = 0
		for i, b := range bytes {
			extension += uint64(b) << (56 - i*8)
		}
		extended = 2
	}

	finallen := uint64(payloadlen) + extension

	mask, err := read(r, 4)
	if err != nil {
		return nil, err
	}

	var data []byte
	if finallen > math.MaxUint32 {
		data = []byte{}
		remaining := finallen

		var i uint64
		for i = 0; i < finallen; i++ {
			next := int(min(remaining, 64))
			b, err := read(r, next)
			if err != nil {
				return nil, err
			}

			bytes = append(bytes, b...)
		}
	} else {
		data, err = read(r, int(finallen))
		if err != nil {
			return nil, err
		}
	}

	for i, d := range data {
		j := i % 4
		data[i] = d ^ mask[j]
	}

	return &Frame{
		isFinal:                 fin,
		rsv1:                    rsv1,
		rsv2:                    rsv2,
		rsv3:                    rsv3,
		Opcode:                  uint8(opcode),
		isMasked:                masked,
		payloadLengthExtensions: extended,
		payloadLength:           finallen,
		maskingKey:              mask,
		PayloadData:             data,
	}, nil
}

func (f *Frame) WriteToBuffer(b *bufio.Writer) {
	bytes := []byte{}

	isFinal := byte(btoi(f.isFinal)) << 7
	bytes = append(bytes, isFinal+f.Opcode)

	if f.payloadLengthExtensions == 0 {
		bytes = append(bytes, byte(f.payloadLength))
	} else {
		var numBytes int

		if f.payloadLengthExtensions == 2 {
			numBytes = 3
		} else {
			numBytes = 7
		}

		for i := numBytes; i >= 0; i-- {
			b := f.payloadLength >> (8 * i)
			bytes = append(bytes, byte(b))
		}
	}

	bytes = append(bytes, f.PayloadData...)

	b.Write(bytes)
	b.Flush()
}

func btoi(b bool) int {
	if b {
		return 1
	}

	return 0
}
