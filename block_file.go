package slrparser

import (
	"os"
)

// MagicBytes is a block delimeter in block file
type MagicBytes uint32

const (
	// MainnetMagicBytes identifies block start in mainnet
	MainnetMagicBytes MagicBytes = 0xfd04f104
	// TestnetMagicBytes identifies block start in testnet
	TestnetMagicBytes MagicBytes = 0x0709110b
)

// BlockFile represents block file
type BlockFile struct {
	BlockEntityReader
	file *os.File
}

func NewBlockFile(filepath string) (*BlockFile, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &BlockFile{BlockEntityReader{reader: file}, file}, nil
}

func (f *BlockFile) Close() {
	f.file.Close()
}

func (f *BlockFile) Seek(offset int64, whence int) (int64, error) {
	return f.file.Seek(offset, whence)
}

func (f *BlockFile) Pos() int64 {
	pos, _ := f.file.Seek(0, 1)
	return pos
}

func (f *BlockFile) Size() (int64, error) {
	fileInfo, err := f.file.Stat()
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), err
}

func (f *BlockFile) Peek(length int) ([]byte, error) {
	pos := f.Pos()
	val := make([]byte, length)
	f.file.Read(val)
	_, err := f.file.Seek(pos, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}
