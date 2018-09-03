package slrparser

import (
	"bytes"
	"time"
)

type BlockIndexParser struct {
	reader *BlockEntityReader
}

func NewBlockIndexParser(bs []byte) *BlockIndexParser {
	reader := &BlockEntityReader{bytes.NewReader(bs)}
	return &BlockIndexParser{reader}
}

func (p *BlockIndexParser) Parse() (*BlockIndex, error) {
	index := &BlockIndex{}
	index.Version = p.reader.ReadInt32()
	index.HashNext = p.reader.ReadBytes(32)
	index.File = p.reader.ReadUint32()
	index.BlockPos = p.reader.ReadUint32()
	index.Height = p.reader.ReadInt32()
	index.Mint = p.reader.ReadUint64()
	index.MoneySupply = p.reader.ReadUint64()
	index.Flags = p.reader.ReadUint32()
	index.StakeModifier = p.reader.ReadBytes(8)
	if index.IsProofOfStake() {
		index.PrevOutStakeHash = p.reader.ReadBytes(32)
		index.PrevOutStakeIndex = p.reader.ReadUint32()
		index.StakeTime = time.Unix(int64(p.reader.ReadUint32()), 0)
		index.HashProofOfStake = p.reader.ReadBytes(32)
	}
	index.Version = p.reader.ReadInt32()
	index.HashPrev = p.reader.ReadBytes(32)
	index.HashMerkle = p.reader.ReadBytes(32)
	index.Timestamp = time.Unix(int64(p.reader.ReadUint32()), 0)
	index.Bits = p.reader.ReadUint32()
	index.Nonce = p.reader.ReadUint32()
	index.Hash = p.reader.ReadBytes(32)
	return index, nil
}
