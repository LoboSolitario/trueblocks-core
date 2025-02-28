// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package statePkg

// EXISTING_CODE
import (
	"net/http"

	outputHelpers "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output/helpers"
	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunState handles the state command for the command line. Returns error only as per cobra.
func RunState(cmd *cobra.Command, args []string) (err error) {
	opts := stateFinishParse(args)
	outputHelpers.SetWriterForCommand("state", &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.StateInternal()
	return
}

// ServeState handles the state command for the API. Returns error and a bool if handled
func ServeState(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := stateFinishParseApi(w, r)
	outputHelpers.InitJsonWriterApi("state", w, &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, handled = opts.StateInternal()
	outputHelpers.CloseJsonWriterIfNeededApi("state", err, &opts.Globals)
	return
}

// StateInternal handles the internal workings of the state command.  Returns error and a bool if handled
func (opts *StateOptions) StateInternal() (err error, handled bool) {
	err = opts.validateState()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	if opts.Globals.IsApiMode() {
		return nil, false
	}

	handled = true
	err = opts.Globals.PassItOn("getState", opts.Globals.Chain, opts.toCmdLine(), opts.getEnvStr())
	// EXISTING_CODE

	return
}

// EXISTING_CODE
// EXISTING_CODE
