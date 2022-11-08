package cache

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"os"
	"testing"
)

func helperCheckTransaction(t *testing.T, tr *CTransaction, expected *CTransaction) (err error) {
	if string(tr.hash.content) != string(expected.hash.content) {
		t.Fatal("mismatched transaction hash")
	}
	if string(tr.blockHash.content) != string(expected.blockHash.content) {
		t.Fatal("mismatched blockHash")
	}
	if tr.blockNumber != expected.blockNumber {
		t.Fatal("mismatched blockNumber")
	}
	if tr.transactionIndex != expected.transactionIndex {
		t.Fatal("mismatched transactionIndex")
	}
	if tr.nonce != expected.nonce {
		t.Fatal("mismatched nonce")
	}
	if tr.timestamp != expected.timestamp {
		t.Fatal("mismatched timestamp", tr.timestamp, expected.timestamp)
	}
	if string(tr.from.content) != string(expected.from.content) {
		t.Fatal("mismatched from")
	}
	if string(tr.to.content) != string(expected.to.content) {
		t.Fatal("mismatched to")
	}
	if tr.value.len != expected.value.len {
		t.Fatal("mismatched value")
	}
	if tr.gas != expected.gas {
		t.Fatal("mismatched gas")
	}
	if tr.gasPrice != expected.gasPrice {
		t.Fatal("mismatched gasPrice")
	}
	if tr.maxFeePerGas != expected.maxFeePerGas {
		t.Fatal("mismatched maxFeePerGas")
	}
	if tr.maxPriorityFeePerGas != expected.maxPriorityFeePerGas {
		t.Fatal("mismatched maxPriorityFeePerGas")
	}
	if string(tr.input.content) != string(expected.input.content) {
		t.Fatal("mismatched input")
	}
	if tr.hasToken != expected.hasToken {
		t.Fatal("mismatched hasToken")
	}

	receipt := tr.receipt
	if string(receipt.contractAddress.content) != string(expected.receipt.contractAddress.content) {
		t.Fatal("mismatched contractAddress")
	}
	if receipt.gasUsed != expected.receipt.gasUsed {
		t.Fatal("mismatched gasUsed")
	}
	if receipt.effectiveGasPrice != expected.receipt.effectiveGasPrice {
		t.Fatal("mismatched effectiveGasPrice")
	}
	if receipt.status != expected.receipt.status {
		t.Fatal("mismatched status")
	}

	logs := receipt.logs
	for logIndex, log := range logs.items {
		if string(log.address.content) != string(expected.receipt.logs.items[logIndex].address.content) {
			t.Fatal("mismatched address")
		}
		if log.logIndex != expected.receipt.logs.items[logIndex].logIndex {
			t.Fatal("mismatched logIndex")
		}
		for topicIndex, topic := range log.topics.items {
			if string(topic.content) != string(expected.receipt.logs.items[logIndex].topics.items[topicIndex].content) {
				t.Fatal("mismatched topic")
			}
		}

		if string(log.data.content) != string(expected.receipt.logs.items[logIndex].data.content) {
			t.Fatal("mismatched data")
		}
	}

	return
}

func TestReadLegacyBlock(t *testing.T) {
	f, err := os.Open("./cache_block.bin")
	defer f.Close()
	buf := bufio.NewReader(f)
	if err != nil {
		t.Fatal("cannot open file")
	}
	block, err := ReadBlock(buf)
	if err != nil {
		t.Fatal(err)
	}

	if block.gasLimit != 6712355 {
		t.Fatal("invalid gasLimit")
	}
	if block.gasUsed != 6630311 {
		t.Fatal("invalid gasUsed")

	}
	if string(block.hash.content) != "0x51bc754831f33817e755039d90af3b20ea1e21905529ddaa03d7ba9f5fc9e66f" {
		t.Fatal("invalid hash")

	}
	if block.blockNumber != 4369999 {
		t.Fatal("invalid blockNumber")

	}
	if string(block.parentHash.content) != "0x0c8afc74a5e5c7d19f8dd5f2ca13098f871e5b6fd8a6e237df111f935d777105" {
		t.Fatal("invalid parentHash")

	}
	if string(block.miner.content) != "0xea674fdde714fd979de3edf0f56aa9716b898ec8" {
		t.Fatal("invalid miner")

	}
	if block.difficulty != 2997274096101735 {
		t.Fatal("invalid difficulty")

	}
	if block.finalized != true {
		t.Fatal("invalid finalized")

	}
	if block.timestamp != 1508131303 {
		t.Fatal("invalid timestamp")

	}
	if block.baseFeePerGas.len != 0 {
		t.Fatal("invalid baseFeePerGas")
	}

	helperCheckTransaction(t, &block.transactions.items[0], &CTransaction{
		hash:                 cString{content: []byte("0x5a9d77909e77eee0a768d4160155c60fa77ae5b621e3c6f4fe6c54f509ea8c33")},
		blockHash:            cString{content: []byte("0x51bc754831f33817e755039d90af3b20ea1e21905529ddaa03d7ba9f5fc9e66f")},
		blockNumber:          4369999,
		transactionIndex:     0,
		nonce:                175,
		timestamp:            0, // 1508131303,
		from:                 cString{content: []byte("0x53871c23523453988ebd6524fcb0ea29241ca4d2")},
		to:                   cString{content: []byte("0x8d12a197cb00d4747a1fe03395095ce2a5cc6819")},
		value:                biguint_t[uint64]{len: 0},
		gas:                  250000,
		gasPrice:             50000000000,
		maxFeePerGas:         0,
		maxPriorityFeePerGas: 0,
		input:                cString{content: []byte("0x338b5dea000000000000000000000000f3db5fa2c66b7af3eb0c0b782510816cbe4813b800000000000000000000000000000000000000000000000000000000004c4b40")},
		hasToken:             0,
		receipt: CReceipt{
			contractAddress:   cString{content: []byte("0x0")},
			gasUsed:           53673,
			effectiveGasPrice: 50000000000,
			logs: cArray[CLogEntry]{
				items: []CLogEntry{
					CLogEntry{
						address:  cString{content: []byte("0xf3db5fa2c66b7af3eb0c0b782510816cbe4813b8")},
						logIndex: 0,
						topics: cArray[CTopic]{
							items: []CTopic{
								cString{content: []byte("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")},
								cString{content: []byte("0x00000000000000000000000053871c23523453988ebd6524fcb0ea29241ca4d2")},
								cString{content: []byte("0x0000000000000000000000008d12a197cb00d4747a1fe03395095ce2a5cc6819")},
							},
						},
						data: cString{content: []byte("0x00000000000000000000000000000000000000000000000000000000004c4b40")},
					},
					CLogEntry{
						address:  cString{content: []byte("0x8d12a197cb00d4747a1fe03395095ce2a5cc6819")},
						logIndex: 1,
						topics: cArray[CTopic]{
							items: []CTopic{
								cString{content: []byte("0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7")},
							},
						},
						data: cString{content: []byte("0x000000000000000000000000f3db5fa2c66b7af3eb0c0b782510816cbe4813b800000000000000000000000053871c23523453988ebd6524fcb0ea29241ca4d200000000000000000000000000000000000000000000000000000000004c4b4000000000000000000000000000000000000000000000000000000000004c4b40")},
					},
				},
			},
			status: 1,
		},
	})

	helperCheckTransaction(t, &block.transactions.items[8], &CTransaction{
		hash:                 cString{content: []byte("0x5589ddfd9db108fc6be96c68df9fdcc89a0673fde85fb4089653334bb8c1fc71")},
		blockHash:            cString{content: []byte("0x51bc754831f33817e755039d90af3b20ea1e21905529ddaa03d7ba9f5fc9e66f")},
		blockNumber:          4369999,
		transactionIndex:     8,
		nonce:                254,
		timestamp:            0, // 1508131303,
		from:                 cString{content: []byte("0x1366a2ca67594ffd5174d0216d60d9ea8deb511f")},
		to:                   cString{content: []byte("0xb4a68dfdfb56184930c3a84e9244919823c3a2b4")},
		value:                biguint_t[uint64]{len: 1}, // 921920000000000
		gas:                  21000,
		gasPrice:             20000000000,
		maxFeePerGas:         0,
		maxPriorityFeePerGas: 0,
		input:                cString{content: []byte("0x")},
		hasToken:             0,
		receipt: CReceipt{
			contractAddress:   cString{content: []byte("0x0")},
			gasUsed:           21000,
			effectiveGasPrice: 20000000000,
			status:            1,
		},
	})
}

func Test_readCppString(t *testing.T) {
	source := bytes.NewBuffer([]byte{})
	err := binary.Write(source, binary.LittleEndian, uint64(6))
	if err != nil {
		t.Fatal("while preparing test data:", err)
	}
	err = binary.Write(source, binary.LittleEndian, []byte("CBlock"))
	if err != nil {
		t.Fatal("while preparing test data:", err)
	}

	result := cString{}
	buf := bufio.NewReader(source)
	err = readString(buf, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result.size != 6 {
		t.Fatal("wrong size:", result.size)
	}
	if string(result.content) != "CBlock" {
		t.Fatal("wrong content:", result.content)
	}
}
