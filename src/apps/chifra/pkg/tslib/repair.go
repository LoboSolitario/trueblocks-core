package tslib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpcClient"
)

// Repair repairs a single timestamp
func Repair(chain string, bn uint64) error {
	cnt, err := NTimestamps(chain)
	if err != nil {
		return err
	}
	if bn >= cnt {
		msg := fmt.Sprintf("Block number %d out of range %d", bn, cnt)
		return errors.New(msg)
	}

	tsFn := filepath.Join(config.GetPathToIndex(chain), "ts.bin")
	tmpPath := filepath.Join(config.GetPathToCache(chain), "tmp")
	if backupFn, err := file.MakeBackup(tmpPath, tsFn); err == nil {
		defer func() {
			DeCache(chain)
			if file.FileExists(backupFn) {
				// If the backup file exists, something failed, so we replace the original file.
				os.Rename(backupFn, tsFn)
				os.Remove(backupFn) // seems redundant, but may not be on some operating systems
			}
		}()

		if fp, err := os.OpenFile(tsFn, os.O_RDWR|os.O_CREATE, 0644); err == nil {
			defer func() { // defers are last in, first out
				fp.Close()
				// sigintTrap.Disable(trapCh)
				// writeMutex.Unlock()
			}()

			recordSize := int64(unsafe.Sizeof(uint32(0))) * 2
			pos := (recordSize * int64(bn))
			fp.Seek(pos, io.SeekStart)

			block, _ := rpcClient.GetBlockByNumber(chain, bn)
			record := Timestamp{Bn: uint32(block.BlockNumber), Ts: uint32(block.TimeStamp)}
			err = binary.Write(fp, binary.LittleEndian, &record)
			if err != nil {
				return err
			}
			fp.Sync() // probably redundant

			os.Remove(backupFn)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
