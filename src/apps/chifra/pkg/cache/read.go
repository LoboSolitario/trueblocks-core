package cache

import (
	"bufio"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type readBytes = func(data any) error

func createReadFn(reader *bufio.Reader) readBytes {
	return func(data any) error {
		return binary.Read(reader, binary.LittleEndian, data)
	}
}

type ArrayItem interface {
	~string |
		common.Hash |
		common.Address |
		types.SimpleTransaction |
		types.SimpleTrace |
		types.SimpleParameter |
		types.SimpleLog
}

func mapArray[Item ArrayItem](
	reader *bufio.Reader,
	target *[]Item,
	readValue func(reader *bufio.Reader) (*Item, error),
) (err error) {
	var itemCount uint64 = 0
	read := createReadFn(reader)
	err = read(&itemCount)
	if err != nil {
		return
	}

	*target = make([]Item, 0, itemCount)

	// result := make([]Item, itemCount)
	for i := 0; uint64(i) < itemCount; i++ {
		item, readErr := readValue(reader)
		if readErr != nil {
			err = readErr
			return
		}

		*target = append(*target, *item)
	}

	// target = target[:len(result)]
	// *target = append(*target, result...)
	// copy(*target, result)
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

func readAddress(reader *bufio.Reader, target *common.Address) (err error) {
	str := &cString{}
	err = readString(reader, str)
	if err != nil {
		return
	}
	addr := common.HexToAddress(string(str.content))
	*target = addr
	return
}

func readHash(reader *bufio.Reader, target *common.Hash) (err error) {
	str := &cString{}
	err = readString(reader, str)
	if err != nil {
		return
	}
	hash := common.HexToHash(string(str.content))
	*target = hash
	return
}

// TODO: rename
func readJustString(reader *bufio.Reader, target *string) (err error) {
	str := &cString{}
	err = readString(reader, str)
	if err != nil {
		return
	}
	s := string(str.content)
	*target = s
	return
}

type bignumBase interface {
	uint64
}
type biguint_t[Base bignumBase] struct {
	capacity int32
	len      int32
	items    []Base
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

func parseBigUint(reader *bufio.Reader, target *big.Int) (err error) {
	cBigUint, err := readBigUint(reader)
	if err != nil {
		return
	}
	if cBigUint.len == 0 {
		*target = *big.NewInt(0)
		return
	}
	// TODO: ask @jrush if this is correct
	*target = *big.NewInt(int64(cBigUint.items[0]))
	return
}

func turnToArrayItem[Item ArrayItem](
	reader *bufio.Reader,
	readAndSetValue func(reader *bufio.Reader, target *Item) (err error),
) func(reader *bufio.Reader) (value *Item, err error) {
	return func(reader *bufio.Reader) (value *Item, err error) {
		v := new(Item)
		err = readAndSetValue(reader, v)
		value = v
		return
	}
}

func readTimestamp(read readBytes, target *time.Time) (err error) {
	var rawTimestamp int64
	err = read(&rawTimestamp)
	if err != nil {
		return
	}
	// TODO: this seems not to work
	*target = time.Unix(rawTimestamp, 0)

	return
}

func readUintAsBool(read readBytes, target *bool) (err error) {
	var raw uint8
	err = read(&raw)
	if err != nil {
		return
	}
	*target = raw != 0
	return
}

func readCacheHeader(reader *bufio.Reader, target *cacheHeader) (err error) {
	read := createReadFn(reader)

	err = read(&target.deleted)
	if err != nil {
		return
	}

	err = read(&target.schema)
	if err != nil {
		return
	}

	err = read(&target.showing)
	if err != nil {
		return
	}

	err = readString(reader, &target.className)
	if err != nil {
		return err
	}

	return
}

func ReadBlockSimple(reader *bufio.Reader) (block *types.SimpleBlock, err error) {
	block = &types.SimpleBlock{}
	read := createReadFn(reader)

	readCacheHeader(reader, &cacheHeader{})

	err = read(&block.GasLimit)
	if err != nil {
		return
	}

	err = read(&block.GasUsed)
	if err != nil {
		return
	}

	err = readHash(reader, &block.Hash)
	if err != nil {
		return
	}

	err = read(&block.BlockNumber)
	if err != nil {
		return
	}

	err = readHash(reader, &block.ParentHash)
	if err != nil {
		return
	}

	err = readAddress(reader, &block.Miner)
	if err != nil {
		return
	}

	err = read(&block.Difficulty)
	if err != nil {
		return
	}

	err = read(&block.Finalized)
	if err != nil {
		return
	}

	err = readTimestamp(read, &block.Timestamp)
	if err != nil {
		return
	}

	err = parseBigUint(reader, &block.BaseFeePerGas)
	if err != nil {
		return
	}

	err = mapArray(reader, &block.Transactions, ReadSimpleTransaction)
	if err != nil {
		return
	}

	return
}

func ReadSimpleTransaction(reader *bufio.Reader) (tx *types.SimpleTransaction, err error) {
	tx = &types.SimpleTransaction{}
	read := createReadFn(reader)

	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readHash(reader, &tx.Hash)
	if err != nil {
		return
	}

	err = readHash(reader, &tx.BlockHash)
	if err != nil {
		return
	}

	err = read(&tx.BlockNumber)
	if err != nil {
		return
	}

	err = read(&tx.TransactionIndex)
	if err != nil {
		return
	}

	err = read(&tx.Nonce)
	if err != nil {
		return
	}

	err = readTimestamp(read, &tx.Timestamp)
	if err != nil {
		return
	}

	err = readAddress(reader, &tx.From)
	if err != nil {
		return
	}

	err = readAddress(reader, &tx.To)
	if err != nil {
		return
	}

	err = parseBigUint(reader, &tx.Value)
	if err != nil {
		return
	}

	err = parseBigUint(reader, &tx.ExtraValue1)
	if err != nil {
		return
	}

	err = parseBigUint(reader, &tx.ExtraValue2)
	if err != nil {
		return
	}

	err = read(&tx.Gas)
	if err != nil {
		return
	}

	err = read(&tx.GasPrice)
	if err != nil {
		return
	}

	err = read(&tx.MaxFeePerGas)
	if err != nil {
		return
	}

	err = read(&tx.MaxPriorityFeePerGas)
	if err != nil {
		return
	}

	err = readJustString(reader, &tx.Input)
	if err != nil {
		return
	}

	err = readUintAsBool(read, &tx.IsError)
	if err != nil {
		return
	}

	err = readUintAsBool(read, &tx.HasToken)
	if err != nil {
		return
	}

	err = read(&tx.Cachebits)
	if err != nil {
		return
	}

	err = read(&tx.Reserved2)
	if err != nil {
		return
	}

	receipt, err := ReadSimpleReceipt(reader)
	if err != nil {
		return
	}
	tx.Receipt = receipt

	err = mapArray(reader, &tx.Traces, ReadSimpleTrace)
	if err != nil {
		return
	}

	articulatedTx, err := ReadSimpleFunction(reader)
	if err != nil {
		return
	}
	tx.ArticulatedTx = articulatedTx

	return
}

func ReadSimpleReceipt(reader *bufio.Reader) (receipt *types.SimpleReceipt, err error) {
	receipt = &types.SimpleReceipt{}
	read := createReadFn(reader)

	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readAddress(reader, &receipt.ContractAddress)
	if err != nil {
		return
	}

	err = read(&receipt.GasUsed)
	if err != nil {
		return
	}

	err = read(&receipt.EffectiveGasPrice)
	if err != nil {
		return
	}

	err = mapArray(reader, &receipt.Logs, ReadSimpleLog)
	if err != nil {
		return
	}

	err = read(&receipt.Status)
	if err != nil {
		return
	}

	return
}

func ReadSimpleLog(reader *bufio.Reader) (log *types.SimpleLog, err error) {
	log = &types.SimpleLog{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readAddress(reader, &log.Address)
	if err != nil {
		return
	}

	err = read(&log.LogIndex)
	if err != nil {
		return
	}

	err = mapArray(reader, &log.Topics, turnToArrayItem(reader, readHash))
	if err != nil {
		return
	}

	err = readJustString(reader, &log.Data)
	if err != nil {
		return
	}

	articulatedLog, err := ReadSimpleFunction(reader)
	if err != nil {
		return
	}
	log.ArticulatedLog = articulatedLog

	return
}

func ReadSimpleFunction(reader *bufio.Reader) (function *types.SimpleFunction, err error) {
	function = &types.SimpleFunction{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readJustString(reader, &function.Name)
	if err != nil {
		return
	}

	err = readJustString(reader, &function.FunctionType)
	if err != nil {
		return
	}

	err = readJustString(reader, &function.AbiSource)
	if err != nil {
		return
	}

	err = read(&function.Anonymous)
	if err != nil {
		return
	}

	err = read(&function.Constant)
	if err != nil {
		return
	}

	err = readJustString(reader, &function.StateMutability)
	if err != nil {
		return
	}

	err = readJustString(reader, &function.Signature)
	if err != nil {
		return
	}

	err = readJustString(reader, &function.Encoding)
	if err != nil {
		return
	}

	err = mapArray(reader, &function.Inputs, ReadSimpleParameter)
	if err != nil {
		return
	}

	err = mapArray(reader, &function.Outputs, ReadSimpleParameter)
	if err != nil {
		return
	}

	return
}

func ReadSimpleParameter(reader *bufio.Reader) (param *types.SimpleParameter, err error) {
	param = &types.SimpleParameter{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readJustString(reader, &param.ParameterType)
	if err != nil {
		return
	}

	err = readJustString(reader, &param.Name)
	if err != nil {
		return
	}

	err = readJustString(reader, &param.StrDefault)
	if err != nil {
		return
	}

	err = readJustString(reader, &param.Value)
	if err != nil {
		return
	}

	err = read(&param.Indexed)
	if err != nil {
		return
	}

	err = readJustString(reader, &param.InternalType)
	if err != nil {
		return
	}

	err = mapArray(reader, &param.Components, ReadSimpleParameter)
	if err != nil {
		return
	}

	err = read(&param.Unused)
	if err != nil {
		return
	}

	err = read(&param.IsFlags)
	if err != nil {
		return
	}

	return
}

func ReadSimpleTrace(reader *bufio.Reader) (trace *types.SimpleTrace, err error) {
	trace = &types.SimpleTrace{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readHash(reader, &trace.BlockHash)
	if err != nil {
		return
	}

	err = read(&trace.BlockNumber)
	if err != nil {
		return
	}

	err = read(&trace.Subtraces)
	if err != nil {
		return
	}

	err = mapArray(reader, &trace.TraceAddress, turnToArrayItem(reader, readJustString))
	if err != nil {
		return
	}

	err = readHash(reader, &trace.TransactionHash)
	if err != nil {
		return
	}

	err = read(&trace.TransactionIndex)
	if err != nil {
		return
	}

	err = readJustString(reader, &trace.TraceType)
	if err != nil {
		return
	}

	err = readJustString(reader, &trace.Error)
	if err != nil {
		return
	}

	action, err := readSimpleTraceAction(reader)
	if err != nil {
		return
	}
	trace.Action = action

	result, err := readSimpleTraceResult(reader)
	if err != nil {
		return
	}
	trace.Result = result

	function, err := ReadSimpleFunction(reader)
	if err != nil {
		return
	}
	trace.ArticulatedTrace = function

	return
}

func readSimpleTraceAction(reader *bufio.Reader) (action *types.SimpleTraceAction, err error) {
	action = &types.SimpleTraceAction{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readAddress(reader, &action.SelfDestructed)
	if err != nil {
		return
	}

	err = parseBigUint(reader, &action.Balance)
	if err != nil {
		return
	}

	err = readJustString(reader, &action.CallType)
	if err != nil {
		return
	}

	err = readAddress(reader, &action.From)
	if err != nil {
		return
	}

	err = read(&action.Gas)
	if err != nil {
		return
	}

	err = readJustString(reader, &action.Init)
	if err != nil {
		return
	}

	err = readJustString(reader, &action.Input)
	if err != nil {
		return
	}

	err = readAddress(reader, &action.RefundAddress)
	if err != nil {
		return
	}

	err = readAddress(reader, &action.To)
	if err != nil {
		return
	}

	err = parseBigUint(reader, &action.Value)
	if err != nil {
		return
	}

	return
}

func readSimpleTraceResult(reader *bufio.Reader) (result *types.SimpleTraceResult, err error) {
	result = &types.SimpleTraceResult{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &cacheHeader{})
	if err != nil {
		return
	}

	err = readAddress(reader, &result.NewContract)
	if err != nil {
		return
	}

	err = readJustString(reader, &result.Code)
	if err != nil {
		return
	}

	err = read(&result.GasUsed)
	if err != nil {
		return
	}

	err = readJustString(reader, &result.Output)
	if err != nil {
		return
	}

	return
}
