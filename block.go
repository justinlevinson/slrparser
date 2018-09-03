package slrparser

import (
	"encoding/binary"
	"encoding/hex"
	"time"
)

type Block struct {
	BlockHeader
	Length       uint32
	Transactions []Transaction
	Signature    BlockSignature
	StartPos     int64
}

type BlockHeader struct {
	Hash       Hash256
	Version    int32
	HashPrev   Hash256
	HashMerkle Hash256
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

// BlockSignature - signed by one of the coin base txout[N]'s owner
type BlockSignature []byte

func (sig BlockSignature) String() string {
	return hex.EncodeToString(sig)
}

func (sig BlockSignature) MarshalText() ([]byte, error) {
	return []byte(sig.String()), nil
}

func (h *BlockHeader) Finalize() {
	version := make([]byte, 4)
	binary.LittleEndian.PutUint32(version, uint32(h.Version))

	timestamp := make([]byte, 4)
	binary.LittleEndian.PutUint32(timestamp, uint32(h.Timestamp.Unix()))

	bits := make([]byte, 4)
	binary.LittleEndian.PutUint32(bits, h.Bits)

	nonce := make([]byte, 4)
	binary.LittleEndian.PutUint32(nonce, h.Nonce)

	bin := make([]byte, 0)
	bin = append(bin, version...)
	bin = append(bin, h.HashPrev...)
	bin = append(bin, h.HashMerkle...)
	bin = append(bin, timestamp...)
	bin = append(bin, bits...)
	bin = append(bin, nonce...)

	h.Hash = DoubleSha256(bin)
}
