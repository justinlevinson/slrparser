package slrparser

import (
	"errors"
	"time"
)

type BlockParser struct {
	file *BlockFile
	mgk  MagicBytes
}

func NewBlockParser(file *BlockFile, mgk MagicBytes) *BlockParser {
	return &BlockParser{file, mgk}
}

func (p *BlockParser) ParseBlock() (*Block, error) {
	startPos := p.file.Pos()
	mgk := p.file.ReadMagicBytes()
	if mgk != p.mgk {
		p.file.Seek(startPos, 0)
		return nil, errors.New("Can't find magic bytes at position " + string(startPos))
	}

	block := &Block{StartPos: startPos}
	if err := p.parseBlockHeader(block); err != nil {
		return nil, err
	}
	if err := p.parseTransactions(block); err != nil {
		return nil, err
	}
	if block.Version >= 3 {
		len := p.file.ReadVarint()
		block.Signature = p.file.ReadBytes(len)
	}
	block.Finalize()
	return block, nil
}

func (p *BlockParser) parseBlockHeader(block *Block) error {
	block.Length = p.file.ReadUint32()
	block.Version = p.file.ReadInt32()
	block.HashPrev = p.file.ReadBytes(32)
	block.HashMerkle = p.file.ReadBytes(32)
	block.Timestamp = time.Unix(int64(p.file.ReadUint32()), 0)
	block.Bits = p.file.ReadUint32()
	block.Nonce = p.file.ReadUint32()

	return nil
}

func (p *BlockParser) parseTransactions(block *Block) error {
	txCount := uint32(p.file.ReadVarint())
	for t := uint32(0); t < txCount; t++ {
		tx, err := p.parseTransaction()
		if err != nil {
			return err
		}
		block.Transactions = append(block.Transactions, *tx)
	}
	return nil
}

func (p *BlockParser) parseTransaction() (*Transaction, error) {
	ver := int32(p.file.ReadUint32())
	txTime := time.Unix(int64(p.file.ReadUint32()), 0)
	tx := &Transaction{
		Version: ver,
		Time:    txTime,
	}

	insCount := int(p.file.ReadVarint())
	for j := 0; j < insCount; j++ {
		hash := p.file.ReadBytes(32)
		idx := p.file.ReadUint32()
		scrLen := p.file.ReadVarint()
		scr := p.file.ReadBytes(scrLen)
		seq := p.file.ReadUint32()
		in := &TransactionInput{
			Hash:     hash,
			Index:    idx,
			Script:   scr,
			Sequence: seq,
		}
		tx.Vin = append(tx.Vin, in)
	}

	outsCount := int(p.file.ReadVarint())
	for j := 0; j < outsCount; j++ {
		val := p.file.ReadUint64()
		scrLen := p.file.ReadVarint()
		scr := p.file.ReadBytes(scrLen)
		out := &TransactionOutput{
			Value:  uint64(val),
			Script: scr,
		}
		tx.Vout = append(tx.Vout, out)
	}

	tx.LockTime = time.Unix(int64(p.file.ReadUint32()), 0)

	if tx.Version >= 2 {
		len := p.file.ReadVarint()
		tx.Comment = string(p.file.ReadBytes(len))
	}

	return tx, nil
}
