// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package abisPkg

// EXISTING_CODE
import (
	"net/http"

	outputHelpers "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output/helpers"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunAbis handles the abis command for the command line. Returns error only as per cobra.
func RunAbis(cmd *cobra.Command, args []string) (err error) {
	opts := abisFinishParse(args)
	outputHelpers.SetWriterForCommand("abis", &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.AbisInternal()
	return
}

// ServeAbis handles the abis command for the API. Returns error and a bool if handled
func ServeAbis(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := abisFinishParseApi(w, r)
	outputHelpers.InitJsonWriterApi("abis", w, &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, handled = opts.AbisInternal()
	outputHelpers.CloseJsonWriterIfNeededApi("abis", err, &opts.Globals)
	return
}

// AbisInternal handles the internal workings of the abis command.  Returns error and a bool if handled
func (opts *AbisOptions) AbisInternal() (err error, handled bool) {
	err = opts.validateAbis()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	if len(opts.Find) > 0 {
		return opts.HandleAbiFind(), true
	}

	if opts.Globals.IsApiMode() {
		return nil, false
	}

	handled = true
	err = opts.Globals.PassItOn("grabABI", opts.Globals.Chain, opts.toCmdLine(), opts.getEnvStr())
	// EXISTING_CODE

	return
}

// EXISTING_CODE
// EXISTING_CODE
