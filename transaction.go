package slrparser

import (
	"encoding/binary"
	"encoding/hex"
	"time"
)

type Transaction struct {
	hash     Hash256
	Version  int32
	Time     time.Time
	Vin      []*TransactionInput
	Vout     []*TransactionOutput
	LockTime time.Time
	Comment  string // Version 2
	StartPos uint64
}

type TransactionInput struct {
	Hash     Hash256
	Index    uint32 // FIXME: ????
	Script   Script
	Sequence uint32
}

type TransactionOutput struct {
	Value  uint64
	Script Script
}

type Script []byte

func (script Script) String() string {
	return hex.EncodeToString(script)
}

func (script Script) MarshalText() ([]byte, error) {
	return []byte(script.String()), nil
}

func (tx Transaction) Hash() Hash256 {
	if tx.hash != nil {
		return tx.hash
	}

	bin := make([]byte, 0)

	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(tx.Version))
	bin = append(bin, version...)

	vinLength := Varint(uint64(len(tx.Vin)))
	bin = append(bin, vinLength...)
	for _, in := range tx.Vin {
		bin = append(bin, in.Binary()...)
	}

	voutLength := Varint(uint64(len(tx.Vout)))
	bin = append(bin, voutLength...)
	for _, out := range tx.Vout {
		bin = append(bin, out.Binary()...)
	}

	locktime := make([]byte, 4)
	binary.LittleEndian.PutUint32(locktime, uint32(tx.LockTime.Unix()))
	bin = append(bin, locktime...)

	tx.hash = DoubleSha256(bin)
	return tx.hash
}

func (in TransactionInput) Binary() []byte {
	index := make([]byte, 4)
	binary.LittleEndian.PutUint32(index, uint32(in.Index))

	scriptLength := Varint(uint64(len(in.Script)))

	sequence := make([]byte, 4)
	binary.LittleEndian.PutUint32(sequence, uint32(in.Sequence))

	bin := make([]byte, 0)
	bin = append(bin, in.Hash...)
	bin = append(bin, index...)
	bin = append(bin, scriptLength...)
	bin = append(bin, in.Script...)
	bin = append(bin, sequence...)

	return bin
}

func (out TransactionOutput) Binary() []byte {
	value := make([]byte, 8)
	binary.LittleEndian.PutUint64(value, uint64(out.Value))

	scriptLength := Varint(uint64(len(out.Script)))

	bin := make([]byte, 0)
	bin = append(bin, value...)
	bin = append(bin, scriptLength...)
	bin = append(bin, out.Script...)

	return bin
}
