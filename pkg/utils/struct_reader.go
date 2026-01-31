package utils

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrOverrun = errors.New("null-terminated string exceeded maxLen")

type StructReader struct {
	reader io.Reader
}

func NewStructReader(reader io.Reader) *StructReader {
	return &StructReader{
		reader: reader,
	}
}

// basic readers

func (r *StructReader) ReadUint8() (uint8, error) {
	var v uint8
	err := binary.Read(r.reader, binary.LittleEndian, &v)
	return v, err
}

func (r *StructReader) ReadUint16() (uint16, error) {
	var v uint16
	err := binary.Read(r.reader, binary.LittleEndian, &v)
	return v, err
}

func (r *StructReader) ReadUint32() (uint32, error) {
	var v uint32
	err := binary.Read(r.reader, binary.LittleEndian, &v)
	return v, err
}

func (r *StructReader) ReadUint64() (uint64, error) {
	var v uint64
	err := binary.Read(r.reader, binary.LittleEndian, &v)
	return v, err
}

func (r *StructReader) ReadBytes(length int) ([]byte, error) {
	v := make([]byte, length)
	_, err := io.ReadFull(r.reader, v)
	return v, err
}

// More advanced ones

func (r *StructReader) ReadUint8Bool() (bool, error) {
	b, err := r.ReadUint8()
	if err != nil {
		return false, err
	}
	return b != 0, nil
}

func (r *StructReader) ReadNullTerminatedString(maxLen int) (string, error) {
	// Read bytes until '\0' or EOF, but not exceeding maxLen.
	var buf []byte
	for i := 0; i < maxLen; i++ {
		b, err := r.ReadUint8()
		if err != nil {
			if err == io.EOF {
				// truncated without null, return what we have
				return string(buf), io.ErrUnexpectedEOF
			}
			return "", err
		}
		if b == 0 {
			return string(buf), nil
		}
		buf = append(buf, b)
	}
	return "", ErrOverrun
}
