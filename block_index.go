package slrparser

import (
	"encoding/hex"
	"time"
)

const (
	blockProofOfStake = 1 << iota
	blockStakeEntropy
	blockStakeModifier
)

type BlockIndex struct {
	Version           int32
	HashNext          Hash256
	File              uint32
	BlockPos          uint32
	Height            int32
	Mint              uint64
	MoneySupply       uint64
	Flags             uint32
	StakeModifier     StakeModifier
	PrevOutStakeHash  Hash256
	PrevOutStakeIndex uint32
	StakeTime         time.Time
	HashProofOfStake  Hash256
	BlockHeader
}

type StakeModifier []byte

func (mod StakeModifier) String() string {
	return hex.EncodeToString(mod)
}

func (mod StakeModifier) MarshalText() ([]byte, error) {
	return []byte(mod.String()), nil
}

func (i *BlockIndex) IsProofOfStake() bool {
	return (i.Flags & blockProofOfStake) != 0
}
