package slrparser

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"os"
	"runtime"
)

func SolarCoinDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("APPDATA")
		return home + "/Bitcoin"
	} else if runtime.GOOS == "osx" || runtime.GOOS == "darwin" {
		return os.Getenv("HOME") + "/Library/Application Support/SolarCoin"
	}
	return os.Getenv("HOME") + "/.solarcoin"
}

func ReverseHex(b []byte) []byte {
	newb := make([]byte, len(b))
	copy(newb, b)
	for i := len(newb)/2 - 1; i >= 0; i-- {
		opp := len(newb) - 1 - i
		newb[i], newb[opp] = newb[opp], newb[i]
	}
	return newb
}

func Varint(n uint64) []byte {
	if n > 4294967295 {
		val := make([]byte, 8)
		binary.BigEndian.PutUint64(val, n)
		return append([]byte{0xFF}, val...)
	} else if n > 65535 {
		val := make([]byte, 4)
		binary.BigEndian.PutUint32(val, uint32(n))
		return append([]byte{0xFE}, val...)
	} else if n > 255 {
		val := make([]byte, 2)
		binary.BigEndian.PutUint16(val, uint16(n))
		return append([]byte{0xFD}, val...)
	} else {
		return []byte{byte(n)}
	}
}

type Hash256 []byte

func (hash Hash256) String() string {
	return hex.EncodeToString(ReverseHex(hash))
}

func (hash Hash256) MarshalText() ([]byte, error) {
	return []byte(hash.String()), nil
}

func DoubleSha256(data []byte) Hash256 {
	hash := sha256.New()
	hash.Write(data)
	firstSha256 := hash.Sum(nil)
	hash.Reset()
	hash.Write(firstSha256)
	return hash.Sum(nil)
}
