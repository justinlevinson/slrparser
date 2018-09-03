package slrparser

import (
	"bytes"
	"encoding/binary"
	"io"
)

type BlockEntityReader struct {
	reader io.Reader
}

func (r *BlockEntityReader) Read(length int) []byte {
	val := make([]byte, length)
	r.reader.Read(val)
	return val
}

func (r *BlockEntityReader) ReadMagicBytes() MagicBytes {
	return MagicBytes(r.ReadUint32())
}

func (r *BlockEntityReader) ReadByte() byte {
	val := make([]byte, 1)
	r.reader.Read(val)
	return val[0]
}

func (r *BlockEntityReader) ReadBytes(length uint64) []byte {
	val := make([]byte, length)
	r.reader.Read(val)
	return val
}

func (r *BlockEntityReader) ReadUint16() uint16 {
	val := make([]byte, 2)
	r.reader.Read(val)
	return binary.LittleEndian.Uint16(val)
}

func (r *BlockEntityReader) ReadInt32() int32 {
	raw := make([]byte, 4)
	r.reader.Read(raw)
	var val int32
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (r *BlockEntityReader) ReadUint32() uint32 {
	val := make([]byte, 4)
	r.reader.Read(val)
	return binary.LittleEndian.Uint32(val)
}

func (r *BlockEntityReader) ReadInt64() int64 {
	raw := make([]byte, 8)
	r.reader.Read(raw)
	var val int64
	binary.Read(bytes.NewReader(raw), binary.LittleEndian, &val)
	return val
}

func (r *BlockEntityReader) ReadUint64() uint64 {
	val := make([]byte, 8)
	r.reader.Read(val)
	return binary.LittleEndian.Uint64(val)
}

func (r *BlockEntityReader) ReadVarint() uint64 {
	intType := r.ReadByte()
	if intType == 0xFF {
		return r.ReadUint64()
	} else if intType == 0xFE {
		return uint64(r.ReadUint32())
	} else if intType == 0xFD {
		return uint64(r.ReadUint16())
	}
	return uint64(intType)
}
