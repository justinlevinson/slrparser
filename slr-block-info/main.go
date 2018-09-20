package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/enapter/slrparser"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Stderr.WriteString("No block hash provided. Usage: slr-block-info [HASH]\n")
		os.Exit(1)
	}

	hashStr := os.Args[1]
	hash, err := hex.DecodeString(hashStr)
	if err != nil {
		panic(err)
	}

	file, err := slrparser.NewBlockFile(slrparser.SolarCoinDir() + "/blk0001.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := leveldb.OpenFile(slrparser.SolarCoinDir()+"/txleveldb", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	val, err := db.Get(append([]byte("b"), slrparser.ReverseHex(hash)...), nil)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Unable to find block %v in block index.\n", hashStr))
	}

	prs := slrparser.NewBlockIndexParser(val)
	index, _ := prs.Parse()

	parser := slrparser.NewBlockParser(file, slrparser.MainnetMagicBytes)
	file.Seek(int64(index.BlockPos)-8, 0)
	block, _ := parser.ParseBlock()

	info := &BlockInfo{block, index}

	j, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(j))
}

type BlockInfo struct {
	*slrparser.Block
	*slrparser.BlockIndex
}
