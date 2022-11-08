package cache

import (
	"bufio"
	"encoding/binary"
	"math/big"
)

type uint8_t = uint8
type uint32_t = uint32

type LegacyReceipt struct {
	CacheHeader

	contractAddress address_t
	// cumulativeGasUsed wei_t
	gasUsed           gas_t
	effectiveGasPrice gas_t
	logs              cArray[CLogEntry]
	// root              string_q
	status uint32_t
}
type CLogEntry struct {
	CacheHeader

	address address_t
	// blockHash           hash_t
	// blockNumber         blknum_t
	logIndex       blknum_t
	topics         cArray[CTopic]
	data           string_q
	articulatedLog CFunction
	// compressedLog       string_q
	// transactionHash     hash_t
	// transactionIndex    blknum_t
	// transactionLogIndex blknum_t
	// timestamp           timestamp_t
	// logType             string_q
	// removed             bool
}
type CTopic = string_q
type cTrace struct {
	CacheHeader

	blockHash        hash_t
	blockNumber      blknum_t
	subtraces        uint64_t
	traceAddress     cArray[string_q]
	transactionHash  hash_t
	transactionIndex blknum_t
	traceType        string_q
	error            string_q
	action           cTraceAction
	result           cTraceResult
	articulatedTrace CFunction
	// compressedTrace  string_q
}
type cTraceAction struct {
	CacheHeader

	selfDestructed address_t
	balance        wei_t
	callType       string_q
	from           address_t
	gas            gas_t
	init           string_q
	input          string_q
	refundAddress  address_t
	to             address_t
	value          wei_t
}
type cTraceResult struct {
	CacheHeader

	newContract address_t
	code        string_q
	gasUsed     gas_t
	output      string_q
}
type CFunction struct {
	CacheHeader

	name            string_q
	functionType    string_q
	abi_source      string_q
	anonymous       bool
	constant        bool
	stateMutability string_q
	signature       string_q
	encoding        string_q
	// message         string_q
	inputs  cArray[CParameter]
	outputs cArray[CParameter]
}
type CParameter struct {
	CacheHeader

	parameterType string_q
	name          string_q
	strDefault    string_q
	value         string_q
	indexed       bool
	internalType  string_q
	components    cArray[CParameter]
	unused        bool
	is_flags      uint64_t
	// precision     uint64_t
	// maxWidth      uint64_t
	// doc           uint64_t
	// disp          uint64_t
	// example       string_q
	// description   string_q
}
type CReconciliation struct {
	CacheHeader

	blockNumber         blknum_t
	transactionIndex    blknum_t
	timestamp           timestamp_t
	assetAddr           address_t
	assetSymbol         string_q
	decimals            uint64_t
	prevBlk             blknum_t
	prevBlkBal          bigint_t
	begBal              bigint_t
	endBal              bigint_t
	amountIn            bigint_t
	internalIn          bigint_t
	selfDestructIn      bigint_t
	minerBaseRewardIn   bigint_t
	minerNephewRewardIn bigint_t
	minerTxFeeIn        bigint_t
	minerUncleRewardIn  bigint_t
	prefundIn           bigint_t
	amountOut           bigint_t
	internalOut         bigint_t
	selfDestructOut     bigint_t
	gasCostOut          bigint_t
	reconciliationType  string_q
	spotPrice           double
	priceSource         string_q
}

// TODO: length of float can be wrong
type double = float32

type cComplexType interface {
	// TODO: string is not complex type
	string_q |
		LegacyTransaction |
		LegacyReceipt |
		cTrace |
		CReconciliation |
		CLogEntry |
		CFunction |
		CParameter
}

type cArray[UnderlyingType cComplexType] struct {
	size  unsignedLong
	items []UnderlyingType
}

type CReceipt = LegacyReceipt
type CTraceArray = cArray[cTrace]
type CReconciliationArray = cArray[CReconciliation]

type LegacyTransaction struct {
	CacheHeader

	hash                 hash_t
	blockHash            hash_t
	blockNumber          blknum_t
	transactionIndex     blknum_t
	nonce                uint64_t
	timestamp            timestamp_t
	from                 address_t
	to                   address_t
	value                wei_t
	extraValue1          wei_t
	extraValue2          wei_t
	gas                  gas_t
	gasPrice             gas_t
	maxFeePerGas         gas_t
	maxPriorityFeePerGas gas_t
	input                string_q
	isError              uint8_t
	hasToken             uint8_t
	cachebits            uint8_t
	reserved2            uint8_t
	receipt              CReceipt
	traces               CTraceArray
	articulatedTx        CFunction
	compressedTx         string_q
	statements           CReconciliationArray
	finalized            bool
}

// TODO: length of this uint can be wrong
type unsignedShort = uint32

type bigint_t struct {
	base  unsignedShort
	value []byte // big.Int
}

func readBigInt(reader *bufio.Reader, target bigint_t) (bigint *big.Int, err error) {
	// First, read base
	var base unsignedShort
	base = 0
	err = binary.Read(reader, binary.LittleEndian, &base)
	if err != nil {
		return
	}
	// Now, read the value (it is stored as string)
	stringValue := string_q{}
	err = readCppString(reader, &stringValue)
	if err != nil {
		return
	}

	// Try to parse the value
	// bigint := big.Int{}
	bigint.SetString(string(stringValue.content), int(base))
	return
}

// TODO: this fn should operate only on simple types (like []string_q) and
// it should not accept cComplexType
func readArray[UnderlyingType cComplexType](reader *bufio.Reader, target cArray[UnderlyingType]) (err error) {
	err = binary.Read(reader, binary.LittleEndian, &target.size)
	if err != nil {
		return
	}
	target.items = make([]UnderlyingType, target.size)
	err = binary.Read(reader, binary.LittleEndian, target.items)
	if err != nil {
		return
	}

	return
}

func createReadFn(reader *bufio.Reader) func(data any) error {
	return func(data any) error {
		return binary.Read(reader, binary.LittleEndian, data)
	}
}

func readCacheHeader(reader *bufio.Reader, target *CacheHeader) (err error) {
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

	err = readCppString(reader, &target.className)
	if err != nil {
		return err
	}

	return
}

func ReadTransaction(reader *bufio.Reader) (tr LegacyTransaction, err error) {
	tr = LegacyTransaction{}
	read := func(data any) error {
		return binary.Read(reader, binary.LittleEndian, data)
	}

	err = readCacheHeader(reader, &tr.CacheHeader)
	if err != nil {
		return
	}
	err = readCppString(reader, &tr.hash)

	if err != nil {
		return
	}

	err = readCppString(reader, &tr.blockHash)
	if err != nil {
		return
	}

	err = read(&tr.blockNumber)
	if err != nil {
		return
	}

	err = read(&tr.transactionIndex)
	if err != nil {
		return
	}

	err = read(&tr.nonce)
	if err != nil {
		return
	}

	err = read(&tr.timestamp)
	if err != nil {
		return
	}

	err = readCppString(reader, &tr.from)
	if err != nil {
		return
	}

	err = readCppString(reader, &tr.to)
	if err != nil {
		return
	}

	// err = read(&tr.value)
	value, err := readBigUint(reader)
	if err != nil {
		return
	}
	tr.value = value

	// err = read(&tr.extraValue1)
	extraValue, err := readBigUint(reader)
	if err != nil {
		return
	}
	tr.extraValue1 = extraValue

	// err = read(&tr.extraValue2)
	extraValue, err = readBigUint(reader)
	if err != nil {
		return
	}
	tr.extraValue2 = extraValue

	err = read(&tr.gas)
	if err != nil {
		return
	}

	err = read(&tr.gasPrice)
	if err != nil {
		return
	}

	err = read(&tr.maxFeePerGas)
	if err != nil {
		return
	}

	err = read(&tr.maxPriorityFeePerGas)
	if err != nil {
		return
	}

	err = readCppString(reader, &tr.input)
	if err != nil {
		return
	}

	err = read(&tr.isError)
	if err != nil {
		return
	}

	err = read(&tr.hasToken)
	if err != nil {
		return
	}

	err = read(&tr.cachebits)
	if err != nil {
		return
	}

	err = read(&tr.reserved2)
	if err != nil {
		return
	}

	receipt, err := ReadReceipt(reader)
	if err != nil {
		return
	}
	tr.receipt = *receipt

	traces, err := readComplexArray(reader, ReadCTrace)
	if err != nil {
		return
	}
	tr.traces = traces

	articulatedTx, err := ReadFunction(reader)
	if err != nil {
		return
	}
	tr.articulatedTx = articulatedTx

	// err = readCppString(reader, &tr.compressedTx)
	// if err != nil {
	// 	return
	// }

	// statements, err := readComplexArray(reader, readReconciliation)
	// if err != nil {
	// 	return
	// }
	// tr.statements = statements

	// err = read(&tr.finalized)
	// if err != nil {
	// 	return
	// }

	return
}

func ReadReceipt(reader *bufio.Reader) (receipt *LegacyReceipt, err error) {
	receipt = &LegacyReceipt{}
	read := createReadFn(reader)

	err = readCacheHeader(reader, &receipt.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &receipt.contractAddress)
	if err != nil {
		return
	}

	err = read(&receipt.gasUsed)
	if err != nil {
		return
	}

	err = read(&receipt.effectiveGasPrice)
	if err != nil {
		return
	}

	logs, err := readComplexArray(reader, ReadCLog)
	if err != nil {
		return
	}
	receipt.logs = logs

	// status is 1, but in output null
	err = read(&receipt.status)
	if err != nil {
		return
	}

	return
}

func ReadCLog(reader *bufio.Reader) (log CLogEntry, err error) {
	log = CLogEntry{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &log.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &log.address)
	if err != nil {
		return
	}

	err = read(&log.logIndex)
	if err != nil {
		return
	}

	topics, err := readComplexArray(reader, func(reader *bufio.Reader) (s string_q, err error) {
		err = readCppString(reader, &s)
		return
	})
	if err != nil {
		return
	}
	log.topics = topics

	err = readCppString(reader, &log.data)
	if err != nil {
		return
	}

	articulatedLog, err := ReadFunction(reader)
	if err != nil {
		return
	}
	log.articulatedLog = articulatedLog

	return
}

func ReadFunction(reader *bufio.Reader) (function CFunction, err error) {
	function = CFunction{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &function.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &function.name)
	if err != nil {
		return
	}

	err = readCppString(reader, &function.functionType)
	if err != nil {
		return
	}

	err = readCppString(reader, &function.abi_source)
	if err != nil {
		return
	}

	err = read(&function.anonymous)
	if err != nil {
		return
	}

	err = read(&function.constant)
	if err != nil {
		return
	}

	err = readCppString(reader, &function.stateMutability)
	if err != nil {
		return
	}

	err = readCppString(reader, &function.signature)
	if err != nil {
		return
	}

	err = readCppString(reader, &function.encoding)
	if err != nil {
		return
	}

	inputs, err := readComplexArray(reader, ReadCParameter)
	if err != nil {
		return
	}
	function.inputs = inputs

	outputs, err := readComplexArray(reader, ReadCParameter)
	if err != nil {
		return
	}
	function.outputs = outputs

	return
}

func ReadCParameter(reader *bufio.Reader) (param CParameter, err error) {
	param = CParameter{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &param.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &param.parameterType)
	if err != nil {
		return
	}

	err = readCppString(reader, &param.name)
	if err != nil {
		return
	}

	err = readCppString(reader, &param.strDefault)
	if err != nil {
		return
	}

	err = readCppString(reader, &param.value)
	if err != nil {
		return
	}

	err = read(&param.indexed)
	if err != nil {
		return
	}

	err = readCppString(reader, &param.internalType)
	if err != nil {
		return
	}

	components, err := readComplexArray(reader, ReadCParameter)
	if err != nil {
		return
	}
	param.components = components

	err = read(&param.unused)
	if err != nil {
		return
	}

	err = read(&param.is_flags)
	if err != nil {
		return
	}

	return
}

func ReadCTrace(reader *bufio.Reader) (trace cTrace, err error) {
	trace = cTrace{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &trace.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &trace.blockHash)
	if err != nil {
		return
	}

	err = read(&trace.blockNumber)
	if err != nil {
		return
	}

	err = read(&trace.subtraces)
	if err != nil {
		return
	}

	addrs, err := readComplexArray(reader, func(reader *bufio.Reader) (s string_q, err error) {
		err = readCppString(reader, &s)
		return
	})
	if err != nil {
		return
	}
	trace.traceAddress = addrs

	err = readCppString(reader, &trace.transactionHash)
	if err != nil {
		return
	}

	err = read(&trace.transactionIndex)
	if err != nil {
		return
	}

	err = readCppString(reader, &trace.traceType)
	if err != nil {
		return
	}

	err = readCppString(reader, &trace.error)
	if err != nil {
		return
	}

	action, err := readTraceAction(reader)
	if err != nil {
		return
	}
	trace.action = action

	result, err := readTraceResult(reader)
	if err != nil {
		return
	}
	trace.result = result

	function, err := ReadFunction(reader)
	if err != nil {
		return
	}
	trace.articulatedTrace = function

	return
}

func readTraceAction(reader *bufio.Reader) (action cTraceAction, err error) {
	action = cTraceAction{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &action.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &action.selfDestructed)
	if err != nil {
		return
	}

	// err = read(&action.balance)
	balance, err := readBigUint(reader)
	if err != nil {
		return
	}
	action.balance = balance

	err = readCppString(reader, &action.callType)
	if err != nil {
		return
	}

	err = readCppString(reader, &action.from)
	if err != nil {
		return
	}

	err = read(&action.gas)
	if err != nil {
		return
	}

	err = readCppString(reader, &action.init)
	if err != nil {
		return
	}

	err = readCppString(reader, &action.input)
	if err != nil {
		return
	}

	err = readCppString(reader, &action.refundAddress)
	if err != nil {
		return
	}

	err = readCppString(reader, &action.to)
	if err != nil {
		return
	}

	// err = read(&action.value)
	value, err := readBigUint(reader)
	if err != nil {
		return
	}
	action.value = value

	return
}

func readTraceResult(reader *bufio.Reader) (result cTraceResult, err error) {
	result = cTraceResult{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &result.CacheHeader)
	if err != nil {
		return
	}

	err = readCppString(reader, &result.newContract)
	if err != nil {
		return
	}

	err = readCppString(reader, &result.code)
	if err != nil {
		return
	}

	err = read(&result.gasUsed)
	if err != nil {
		return
	}

	err = readCppString(reader, &result.output)
	if err != nil {
		return
	}

	return
}

func readReconciliation(reader *bufio.Reader) (recon CReconciliation, err error) {
	recon = CReconciliation{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &recon.CacheHeader)
	if err != nil {
		return
	}

	err = read(&recon.blockNumber)
	if err != nil {
		return
	}

	err = read(&recon.transactionIndex)
	if err != nil {
		return
	}

	err = read(&recon.timestamp)
	if err != nil {
		return
	}

	err = readCppString(reader, &recon.assetAddr)
	if err != nil {
		return
	}

	err = readCppString(reader, &recon.assetSymbol)
	if err != nil {
		return
	}

	err = read(&recon.decimals)
	if err != nil {
		return
	}

	err = read(&recon.prevBlk)
	if err != nil {
		return
	}

	err = read(&recon.prevBlkBal)
	if err != nil {
		return
	}

	err = read(&recon.begBal)
	if err != nil {
		return
	}

	err = read(&recon.endBal)
	if err != nil {
		return
	}

	err = read(&recon.amountIn)
	if err != nil {
		return
	}

	err = read(&recon.internalIn)
	if err != nil {
		return
	}

	err = read(&recon.selfDestructIn)
	if err != nil {
		return
	}

	err = read(&recon.minerBaseRewardIn)
	if err != nil {
		return
	}

	err = read(&recon.minerNephewRewardIn)
	if err != nil {
		return
	}

	err = read(&recon.minerTxFeeIn)
	if err != nil {
		return
	}

	err = read(&recon.minerUncleRewardIn)
	if err != nil {
		return
	}

	err = read(&recon.prefundIn)
	if err != nil {
		return
	}

	err = read(&recon.amountOut)
	if err != nil {
		return
	}

	err = read(&recon.internalOut)
	if err != nil {
		return
	}

	err = read(&recon.selfDestructOut)
	if err != nil {
		return
	}

	err = read(&recon.gasCostOut)
	if err != nil {
		return
	}

	err = readCppString(reader, &recon.reconciliationType)
	if err != nil {
		return
	}

	err = read(&recon.spotPrice)
	if err != nil {
		return
	}

	err = readCppString(reader, &recon.priceSource)
	if err != nil {
		return
	}

	return
}

func readComplexArray[T cComplexType](
	reader *bufio.Reader,
	readValue func(reader *bufio.Reader) (T, error),
) (results cArray[T], err error) {
	var itemCount unsignedLong = 0
	err = binary.Read(reader, binary.LittleEndian, &itemCount)
	if err != nil {
		return
	}

	// if itemCount == 0 {
	// 	return
	// }

	items := make([]T, itemCount)
	for i := 0; uint64(i) < itemCount; i++ {
		value, rerr := readValue(reader)
		err = rerr
		if err != nil {
			return
		}
		items[i] = value
	}
	results = cArray[T]{
		size:  itemCount,
		items: items,
	}

	return
}
