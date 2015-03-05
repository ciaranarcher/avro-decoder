package decoder

import (
	"encoding/binary"
	"math"
	"github.com/stealthly/go-avro/avro"
)

var MAX_INT_BUF_SIZE = 5
var MAX_LONG_BUF_SIZE = 10

type BinaryDecoder struct {
	buf []byte
	pos int64
}

func NewBinaryDecoder(buf []byte) *BinaryDecoder {
	return &BinaryDecoder{buf, 0}
}

func (bd *BinaryDecoder) ReadNull() (interface{}, error) {
	return nil, nil
}

func (bd *BinaryDecoder) ReadInt() (int32, error) {
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return 0, avro.EOF
	}
	var value uint32
	var b uint8
	var offset int
	for {
		if offset == MAX_INT_BUF_SIZE {
			return 0, avro.IntOverflow
		}
		b = bd.buf[bd.pos]
		value |= uint32(b & 0x7F)<<uint(7 * offset)
		bd.pos++
		offset++
		if (b&0x80 == 0) {
			break
		}
	}
	return int32((value >> 1) ^ -(value & 1)), nil
}

func (bd *BinaryDecoder) ReadLong() (int64, error) {
	var value uint64
	var b uint8
	var offset int
	for {
		if offset == MAX_LONG_BUF_SIZE {
			return 0, avro.LongOverflow
		}
		b = bd.buf[bd.pos]
		value |= uint64(b & 0x7F)<<uint(7 * offset)
		bd.pos++
		offset++
		if (b&0x80 == 0) {
			break
		}
	}
	return int64((value >> 1) ^ -(value & 1)), nil
}

func (bd *BinaryDecoder) ReadString() (string, error) {
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return "", err
	}
	length, err := bd.ReadInt()
	if err != nil || length < 0 {
		return "", avro.InvalidStringLength
	}
	if err := checkEOF(bd.buf, bd.pos, int(length)); err != nil {
		return "", err
	}
	value := string(bd.buf[bd.pos:int32(bd.pos) + length])
	bd.pos += int64(length)
	return value, nil
}

func (bd *BinaryDecoder) ReadBoolean() (bool, error) {
	b := bd.buf[bd.pos] & 0xFF
	bd.pos++
	var err error = nil
	if b != 0 && b != 1 {
		err = avro.InvalidBool
	}
	return b == 1, err
}

func (bd *BinaryDecoder) ReadBytes() ([]byte, error) {
	//TODO make something with these if's!!
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return nil, avro.EOF
	}
	length, err := bd.ReadLong()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, avro.NegativeBytesLength
	}
	if err := checkEOF(bd.buf, bd.pos, int(length)); err != nil {
		return nil, avro.EOF
	}

	bytes := make([]byte, length)
	copy(bytes[:], bd.buf[bd.pos:bd.pos+length])
	bd.pos += length
	return bytes, err
}

func (bd *BinaryDecoder) ReadFloat() (float32, error) {
	var float float32
	if err := checkEOF(bd.buf, bd.pos, 4); err != nil {
		return float, err
	}
	bits := binary.LittleEndian.Uint32(bd.buf[bd.pos:bd.pos+4])
	float = math.Float32frombits(bits)
	bd.pos += 4
	return float, nil
}

func (bd *BinaryDecoder) ReadDouble() (float64, error) {
	var double float64
	if err := checkEOF(bd.buf, bd.pos, 8); err != nil {
		return double, err
	}
	bits := binary.LittleEndian.Uint64(bd.buf[bd.pos:bd.pos+8])
	double = math.Float64frombits(bits)
	bd.pos += 8
	return double, nil
}

func (bd *BinaryDecoder) ReadEnum() (int32, error) {
	return bd.ReadInt()
}

func (bd *BinaryDecoder) ReadArrayStart() (int64, error) {
	return bd.readItemCount()
}

func (bd *BinaryDecoder) ArrayNext() (int64, error) {
	return bd.readItemCount()
}

func (bd *BinaryDecoder) ReadMapStart() (int64, error) {
	return bd.readItemCount()
}

func (bd *BinaryDecoder) MapNext() (int64, error) {
	return bd.readItemCount()
}

func (bd *BinaryDecoder) readItemCount() (int64, error) {
	if count, err := bd.ReadLong(); err != nil {
		return 0, err
	} else {
		if count < 0 {
			bd.ReadLong()
			count = -count
		}
		return count, err
	}
}

func (bd *BinaryDecoder) ReadFixed(bytes []byte) error {
	return bd.readBytes(bytes, 0, len(bytes))
}

func (bd *BinaryDecoder) ReadFixedWithBounds(bytes []byte, start int, length int) error {
	return bd.readBytes(bytes, start, length)
}

func (bd *BinaryDecoder) readBytes(bytes []byte, start int, length int) error {
	if length < 0 {
		return avro.NegativeBytesLength
	}
	if err := checkEOF(bd.buf, bd.pos, int(start + length)); err != nil {
		return avro.EOF
	}
	copy(bytes[:], bd.buf[bd.pos+int64(start):bd.pos+int64(start)+int64(length)])
	bd.pos += int64(length)

	return nil
}

func (bd *BinaryDecoder) SetBlock(block *avro.DataBlock) {
	bd.buf = block.Data
	bd.Seek(0)
}

func (bd *BinaryDecoder) Seek(pos int64) {
	bd.pos = pos
}

func (bd *BinaryDecoder) Tell() int64 {
	return bd.pos
}

func checkEOF(buf []byte, pos int64, length int) error {
	if int64(len(buf)) < pos+int64(length) {
		return avro.EOF
	}
	return nil
}
