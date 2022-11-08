package cache

import (
	"bufio"
	"encoding/binary"
)

type uint8_t = uint8
type uint32_t = uint32
type double = float32
type unsignedShort = uint32

type CReceipt struct {
	cacheHeader

	contractAddress address_t
	// cumulativeGasUsed wei_t
	gasUsed           gas_t
	effectiveGasPrice gas_t
	logs              cArray[CLogEntry]
	// root              string_q
	status uint32_t
}
type CLogEntry struct {
	cacheHeader

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
type CTrace struct {
	cacheHeader

	blockHash        hash_t
	blockNumber      blknum_t
	subtraces        uint64_t
	traceAddress     cArray[string_q]
	transactionHash  hash_t
	transactionIndex blknum_t
	traceType        string_q
	error            string_q
	action           CTraceAction
	result           CTraceResult
	articulatedTrace CFunction
	// compressedTrace  string_q
}
type cTraceArray = cArray[CTrace]
type CTraceAction struct {
	cacheHeader

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
type CTraceResult struct {
	cacheHeader

	newContract address_t
	code        string_q
	gasUsed     gas_t
	output      string_q
}
type CFunction struct {
	cacheHeader

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
	cacheHeader

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

// type CReconciliation struct {
// 	cacheHeader

// 	blockNumber         blknum_t
// 	transactionIndex    blknum_t
// 	timestamp           timestamp_t
// 	assetAddr           address_t
// 	assetSymbol         string_q
// 	decimals            uint64_t
// 	prevBlk             blknum_t
// 	prevBlkBal          bigint_t
// 	begBal              bigint_t
// 	endBal              bigint_t
// 	amountIn            bigint_t
// 	internalIn          bigint_t
// 	selfDestructIn      bigint_t
// 	minerBaseRewardIn   bigint_t
// 	minerNephewRewardIn bigint_t
// 	minerTxFeeIn        bigint_t
// 	minerUncleRewardIn  bigint_t
// 	prefundIn           bigint_t
// 	amountOut           bigint_t
// 	internalOut         bigint_t
// 	selfDestructOut     bigint_t
// 	gasCostOut          bigint_t
// 	reconciliationType  string_q
// 	spotPrice           double
// 	priceSource         string_q
// }
// type cReconciliationArray = cArray[CReconciliation]

type cArrayItem interface {
	string_q |
		CTransaction |
		CReceipt |
		CTrace |
		// CReconciliation |
		CLogEntry |
		CFunction |
		CParameter
}

type cArray[Item cArrayItem] struct {
	size  unsignedLong
	items []Item
}

type CTransaction struct {
	cacheHeader

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
	traces               cTraceArray
	articulatedTx        CFunction
}

func createReadFn(reader *bufio.Reader) func(data any) error {
	return func(data any) error {
		return binary.Read(reader, binary.LittleEndian, data)
	}
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

func ReadTransaction(reader *bufio.Reader) (tr CTransaction, err error) {
	tr = CTransaction{}
	read := createReadFn(reader)

	err = readCacheHeader(reader, &tr.cacheHeader)
	if err != nil {
		return
	}
	err = readString(reader, &tr.hash)

	if err != nil {
		return
	}

	err = readString(reader, &tr.blockHash)
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

	err = readString(reader, &tr.from)
	if err != nil {
		return
	}

	err = readString(reader, &tr.to)
	if err != nil {
		return
	}

	value, err := readBigUint(reader)
	if err != nil {
		return
	}
	tr.value = value

	extraValue, err := readBigUint(reader)
	if err != nil {
		return
	}
	tr.extraValue1 = extraValue

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

	err = readString(reader, &tr.input)
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

	traces, err := readArray(reader, ReadTrace)
	if err != nil {
		return
	}
	tr.traces = traces

	articulatedTx, err := ReadFunction(reader)
	if err != nil {
		return
	}
	tr.articulatedTx = articulatedTx

	return
}

func ReadReceipt(reader *bufio.Reader) (receipt *CReceipt, err error) {
	receipt = &CReceipt{}
	read := createReadFn(reader)

	err = readCacheHeader(reader, &receipt.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &receipt.contractAddress)
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

	logs, err := readArray(reader, ReadLog)
	if err != nil {
		return
	}
	receipt.logs = logs

	err = read(&receipt.status)
	if err != nil {
		return
	}

	return
}

func ReadLog(reader *bufio.Reader) (log CLogEntry, err error) {
	log = CLogEntry{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &log.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &log.address)
	if err != nil {
		return
	}

	err = read(&log.logIndex)
	if err != nil {
		return
	}

	topics, err := readArray(reader, readStringArrayItem)
	if err != nil {
		return
	}
	log.topics = topics

	err = readString(reader, &log.data)
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
	err = readCacheHeader(reader, &function.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &function.name)
	if err != nil {
		return
	}

	err = readString(reader, &function.functionType)
	if err != nil {
		return
	}

	err = readString(reader, &function.abi_source)
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

	err = readString(reader, &function.stateMutability)
	if err != nil {
		return
	}

	err = readString(reader, &function.signature)
	if err != nil {
		return
	}

	err = readString(reader, &function.encoding)
	if err != nil {
		return
	}

	inputs, err := readArray(reader, ReadParameter)
	if err != nil {
		return
	}
	function.inputs = inputs

	outputs, err := readArray(reader, ReadParameter)
	if err != nil {
		return
	}
	function.outputs = outputs

	return
}

func ReadParameter(reader *bufio.Reader) (param CParameter, err error) {
	param = CParameter{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &param.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &param.parameterType)
	if err != nil {
		return
	}

	err = readString(reader, &param.name)
	if err != nil {
		return
	}

	err = readString(reader, &param.strDefault)
	if err != nil {
		return
	}

	err = readString(reader, &param.value)
	if err != nil {
		return
	}

	err = read(&param.indexed)
	if err != nil {
		return
	}

	err = readString(reader, &param.internalType)
	if err != nil {
		return
	}

	components, err := readArray(reader, ReadParameter)
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

func ReadTrace(reader *bufio.Reader) (trace CTrace, err error) {
	trace = CTrace{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &trace.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &trace.blockHash)
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

	addrs, err := readArray(reader, func(reader *bufio.Reader) (s string_q, err error) {
		err = readString(reader, &s)
		return
	})
	if err != nil {
		return
	}
	trace.traceAddress = addrs

	err = readString(reader, &trace.transactionHash)
	if err != nil {
		return
	}

	err = read(&trace.transactionIndex)
	if err != nil {
		return
	}

	err = readString(reader, &trace.traceType)
	if err != nil {
		return
	}

	err = readString(reader, &trace.error)
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

func readTraceAction(reader *bufio.Reader) (action CTraceAction, err error) {
	action = CTraceAction{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &action.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &action.selfDestructed)
	if err != nil {
		return
	}

	balance, err := readBigUint(reader)
	if err != nil {
		return
	}
	action.balance = balance

	err = readString(reader, &action.callType)
	if err != nil {
		return
	}

	err = readString(reader, &action.from)
	if err != nil {
		return
	}

	err = read(&action.gas)
	if err != nil {
		return
	}

	err = readString(reader, &action.init)
	if err != nil {
		return
	}

	err = readString(reader, &action.input)
	if err != nil {
		return
	}

	err = readString(reader, &action.refundAddress)
	if err != nil {
		return
	}

	err = readString(reader, &action.to)
	if err != nil {
		return
	}

	value, err := readBigUint(reader)
	if err != nil {
		return
	}
	action.value = value

	return
}

func readTraceResult(reader *bufio.Reader) (result CTraceResult, err error) {
	result = CTraceResult{}
	read := createReadFn(reader)
	err = readCacheHeader(reader, &result.cacheHeader)
	if err != nil {
		return
	}

	err = readString(reader, &result.newContract)
	if err != nil {
		return
	}

	err = readString(reader, &result.code)
	if err != nil {
		return
	}

	err = read(&result.gasUsed)
	if err != nil {
		return
	}

	err = readString(reader, &result.output)
	if err != nil {
		return
	}

	return
}

// func readReconciliation(reader *bufio.Reader) (recon CReconciliation, err error) {
// 	recon = CReconciliation{}
// 	read := createReadFn(reader)
// 	err = readCacheHeader(reader, &recon.cacheHeader)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.blockNumber)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.transactionIndex)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.timestamp)
// 	if err != nil {
// 		return
// 	}

// 	err = readString(reader, &recon.assetAddr)
// 	if err != nil {
// 		return
// 	}

// 	err = readString(reader, &recon.assetSymbol)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.decimals)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.prevBlk)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.prevBlkBal)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.begBal)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.endBal)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.amountIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.internalIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.selfDestructIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.minerBaseRewardIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.minerNephewRewardIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.minerTxFeeIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.minerUncleRewardIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.prefundIn)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.amountOut)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.internalOut)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.selfDestructOut)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.gasCostOut)
// 	if err != nil {
// 		return
// 	}

// 	err = readString(reader, &recon.reconciliationType)
// 	if err != nil {
// 		return
// 	}

// 	err = read(&recon.spotPrice)
// 	if err != nil {
// 		return
// 	}

// 	err = readString(reader, &recon.priceSource)
// 	if err != nil {
// 		return
// 	}

// 	return
// }

func readArray[Item cArrayItem](
	reader *bufio.Reader,
	readValue func(reader *bufio.Reader) (Item, error),
) (results cArray[Item], err error) {
	var itemCount unsignedLong = 0
	read := createReadFn(reader)
	err = read(&itemCount)
	if err != nil {
		return
	}

	items := make([]Item, itemCount)
	for i := 0; uint64(i) < itemCount; i++ {
		value, rerr := readValue(reader)
		err = rerr
		if err != nil {
			return
		}
		items[i] = value
	}
	results = cArray[Item]{
		size:  itemCount,
		items: items,
	}

	return
}
