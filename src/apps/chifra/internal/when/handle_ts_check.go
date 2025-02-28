package whenPkg

import (
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/identifiers"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/progress"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpcClient"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/tslib"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

// HandleTimestampsCheck handles chifra when --timestamps --check
func (opts *WhenOptions) HandleTimestampsCheck() error {
	cnt, err := tslib.NTimestamps(opts.Globals.Chain)
	if err != nil {
		return err
	}

	// For display only
	skip := uint64(500)
	if opts.Deep {
		m, _ := rpcClient.GetMetaData(opts.Globals.Chain, opts.Globals.TestMode)
		skip = m.Latest / 500
	}
	scanBar := progress.NewScanBar(cnt /* wanted */, (cnt / skip) /* freq */, cnt /* max */, (2. / 3.))

	blockNums, err := identifiers.GetBlockNumbers(opts.Globals.Chain, opts.BlockIds)
	if err != nil {
		return err
	}

	prev := types.SimpleTimestamp{
		BlockNumber: utils.NOPOS,
		TimeStamp:   utils.NOPOS,
	}

	if len(blockNums) > 0 {
		for _, bn := range blockNums {
			if bn < cnt { // ranges may include blocks after last block
				if err = opts.checkOneBlock(scanBar, &prev, bn); err != nil {
					return err
				}
			}
		}
	} else {
		for bn := uint64(0); bn < cnt; bn++ {
			if err = opts.checkOneBlock(scanBar, &prev, bn); err != nil {
				return err
			}
		}
	}

	return nil
}

func (opts *WhenOptions) checkOneBlock(scanBar *progress.ScanBar, prev *types.SimpleTimestamp, bn uint64) error {
	// The i'th item in the timestamp array on disc
	itemOnDisc, err := tslib.FromBn(opts.Globals.Chain, bn)
	if err != nil {
		return err
	}

	// This just simplifies the code below by removing the need to type cast
	onDisc := types.SimpleTimestamp{
		BlockNumber: uint64(itemOnDisc.Bn),
		TimeStamp:   uint64(itemOnDisc.Ts),
	}

	expected := types.SimpleNamedBlock{BlockNumber: bn, TimeStamp: onDisc.TimeStamp}
	if opts.Deep {
		// If we're going deep, we need to query the node
		expected, _ = rpcClient.GetBlockByNumber(opts.Globals.Chain, bn)
	}

	if prev.TimeStamp != utils.NOPOS {
		status := "Okay"

		bnSequential := prev.BlockNumber < onDisc.BlockNumber
		if !bnSequential {
			msg := fmt.Sprintf("At block %d, block number %d is not one plus %d.%s", bn, onDisc.BlockNumber, prev.BlockNumber, clear)
			logger.Log(logger.Error, msg)
			status = "Error"
		}

		tsSequential := prev.TimeStamp < onDisc.TimeStamp
		if !tsSequential {
			msg := fmt.Sprintf("At block %d, timestamp %d does not increase over previous %d%s", bn, onDisc.TimeStamp, prev.TimeStamp, clear)
			logger.Log(logger.Error, msg)
			status = "Error"
		}

		deepTsCheck := !opts.Deep || (onDisc.TimeStamp == expected.TimeStamp)
		if !deepTsCheck {
			msg := fmt.Sprintf("At block %d, timestamp on disc %d does not agree with on chain %d%s", bn, onDisc.TimeStamp, expected.TimeStamp, clear)
			logger.Log(logger.Error, msg)
			status = "Error"
		}

		posOnDisc := bn == onDisc.BlockNumber
		if opts.Deep {
			posOnDisc = (bn == onDisc.BlockNumber && onDisc.BlockNumber == expected.BlockNumber)
		}
		if !posOnDisc {
			msg := fmt.Sprintf("At block %d, onDisc block number %d does not match expected %d%s", bn, onDisc.BlockNumber, expected.BlockNumber, clear)
			logger.Log(logger.Error, msg)
			status = "Error"
		}

		if status == "Okay" {
			scanBar.Report(opts.Globals.Writer, status, fmt.Sprintf(" bn: %d ts: %d", expected.BlockNumber, expected.TimeStamp))
		}
	}

	*prev = onDisc
	return nil
}

// TODO: There's got to be a better way
var clear = strings.Repeat(" ", 80)
