package cache

import (
	"bufio"
	"encoding/binary"
)

type unsignedLong = uint64
type charPtr = cppString
type gas_t = uint64
type hash_t = cppString
type blknum_t = uint64
type string_q = cppString
type address_t = cppString
type uint64_t = uint64
type timestamp_t = int64

type bignumBase interface {
	uint64
}
type biguint_t[Base bignumBase] struct {
	capacity int32
	len      int32
	items    []Base
}
type wei_t = biguint_t[uint64]

func readBigUint[Base bignumBase](reader *bufio.Reader) (biguint biguint_t[Base], err error) {
	err = binary.Read(reader, binary.LittleEndian, &biguint.capacity)
	if err != nil {
		return
	}
	err = binary.Read(reader, binary.LittleEndian, &biguint.len)
	if err != nil {
		return
	}

	biguint.items = make([]Base, biguint.len)
	err = binary.Read(reader, binary.LittleEndian, &biguint.items)
	return
}

type cStringArray struct {
	size  unsignedLong
	items []cppString
}

type cppString struct {
	size    unsignedLong
	content []byte
}

type cTransactionArray struct {
	size unsignedLong
	// TODO: implement rest of this type
}

type CacheHeader struct {
	deleted   uint64_t
	schema    uint64_t
	showing   uint64_t
	className charPtr
}

type LegacyBlock struct {
	// BaseNode
	// deleted   uint64_t
	// schema    uint64_t
	// showing   uint64_t
	// className charPtr
	CacheHeader

	// Block
	gasLimit      gas_t
	gasUsed       gas_t
	hash          hash_t
	blockNumber   blknum_t
	parentHash    hash_t
	miner         address_t
	difficulty    uint64_t
	finalized     bool
	timestamp     timestamp_t
	baseFeePerGas wei_t
	transactions  cArray[LegacyTransaction]
	// tx_hashes     cStringArray
	// name          string_q
	// unclesCnt     uint64_t
}

func ReadLegacyBlock(reader *bufio.Reader) (block LegacyBlock, err error) {
	read := func(data any) error {
		return binary.Read(reader, binary.LittleEndian, data)
	}
	// block := LegacyBlock{}

	err = read(&block.deleted)
	if err != nil {
		return
	}

	err = read(&block.schema)
	if err != nil {
		return
	}

	err = read(&block.showing)
	if err != nil {
		return
	}

	err = readCppString(reader, &block.className)
	if err != nil {
		return
	}

	err = read(&block.gasLimit)
	if err != nil {
		return
	}

	err = read(&block.gasUsed)
	if err != nil {
		return
	}

	err = readCppString(reader, &block.hash)
	if err != nil {
		return
	}

	err = read(&block.blockNumber)
	if err != nil {
		return
	}

	err = readCppString(reader, &block.parentHash)
	if err != nil {
		return
	}

	err = readCppString(reader, &block.miner)
	if err != nil {
		return
	}

	err = read(&block.difficulty)
	if err != nil {
		return
	}

	err = read(&block.finalized)
	if err != nil {
		return
	}

	err = read(&block.timestamp)
	if err != nil {
		return
	}

	// err = read(&block.baseFeePerGas)
	baseFeePerGas, err := readBigUint(reader)
	if err != nil {
		return
	}
	block.baseFeePerGas = baseFeePerGas

	// var txCount uint64 = 0
	// err = read(&txCount)
	// if err != nil {
	// 	return
	// }

	// txs := make([]*LegacyTransaction, int(txCount))
	// for i := 0; uint64(i) < txCount; i++ {
	// 	tx, terr := ReadTransaction(reader)
	// 	err = terr
	// 	if err != nil {
	// 		return
	// 	}
	// 	txs = append(txs, tx)
	// }
	transactions, err := readComplexArray(reader, ReadTransaction)
	if err != nil {
		return
	}
	block.transactions = transactions

	return
}

func readCppString(reader *bufio.Reader, cppStr *cppString) (err error) {
	err = binary.Read(reader, binary.LittleEndian, &cppStr.size)
	if err != nil {
		return
	}

	cppStr.content = make([]byte, cppStr.size)
	err = binary.Read(reader, binary.LittleEndian, cppStr.content)
	if err != nil {
		return
	}

	return
}
