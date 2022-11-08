package cache

import (
	"bufio"
)

type unsignedLong = uint64
type cString struct {
	size    unsignedLong
	content []byte
}
type charPtr = cString
type gas_t = uint64
type hash_t = cString
type blknum_t = uint64
type string_q = cString
type address_t = cString
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

type cacheHeader struct {
	deleted   uint64_t
	schema    uint64_t
	showing   uint64_t
	className charPtr
}

type CBlock struct {
	// BaseNode
	cacheHeader

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
	transactions  cArray[CTransaction]
}

func readBigUint[Base bignumBase](reader *bufio.Reader) (biguint biguint_t[Base], err error) {
	read := createReadFn(reader)
	err = read(&biguint.capacity)
	if err != nil {
		return
	}
	err = read(&biguint.len)
	if err != nil {
		return
	}

	biguint.items = make([]Base, biguint.len)
	err = read(&biguint.items)
	return
}

func ReadBlock(reader *bufio.Reader) (block CBlock, err error) {
	read := createReadFn(reader)

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

	err = readString(reader, &block.className)
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

	err = readString(reader, &block.hash)
	if err != nil {
		return
	}

	err = read(&block.blockNumber)
	if err != nil {
		return
	}

	err = readString(reader, &block.parentHash)
	if err != nil {
		return
	}

	err = readString(reader, &block.miner)
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

	baseFeePerGas, err := readBigUint(reader)
	if err != nil {
		return
	}
	block.baseFeePerGas = baseFeePerGas

	transactions, err := readArray(reader, ReadTransaction)
	if err != nil {
		return
	}
	block.transactions = transactions

	return
}

// readString reads cString structure from reader. It has different signature than
// the rest of `read*` functions in this package, to ease reading values into
// other structs' fields
func readString(reader *bufio.Reader, str *cString) (err error) {
	read := createReadFn(reader)
	err = read(&str.size)
	if err != nil {
		return
	}

	str.content = make([]byte, str.size)
	err = read(str.content)
	if err != nil {
		return
	}

	return
}

// readStringArrayItem is a wrapper function for `readString`. It does the same thing,
// but can be easily used as parameter to `readArray`
func readStringArrayItem(reader *bufio.Reader) (str cString, err error) {
	err = readString(reader, &str)
	return
}
