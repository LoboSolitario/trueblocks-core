// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were generated with makeClass --run. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package scrapePkg

// EXISTING_CODE
import (
	"net/http"

	outputHelpers "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output/helpers"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/validate"

	"github.com/spf13/cobra"
)

// EXISTING_CODE

// RunScrape handles the scrape command for the command line. Returns error only as per cobra.
func RunScrape(cmd *cobra.Command, args []string) (err error) {
	opts := scrapeFinishParse(args)
	outputHelpers.SetWriterForCommand("scrape", &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, _ = opts.ScrapeInternal()
	return
}

// ServeScrape handles the scrape command for the API. Returns error and a bool if handled
func ServeScrape(w http.ResponseWriter, r *http.Request) (err error, handled bool) {
	opts := scrapeFinishParseApi(w, r)
	outputHelpers.InitJsonWriterApi("scrape", w, &opts.Globals)
	// EXISTING_CODE
	// EXISTING_CODE
	err, handled = opts.ScrapeInternal()
	outputHelpers.CloseJsonWriterIfNeededApi("scrape", err, &opts.Globals)
	return
}

// ScrapeInternal handles the internal workings of the scrape command.  Returns error and a bool if handled
func (opts *ScrapeOptions) ScrapeInternal() (err error, handled bool) {
	err = opts.validateScrape()
	if err != nil {
		return err, true
	}

	// EXISTING_CODE
	if opts.Globals.IsApiMode() {
		return validate.Usage("chifra scrape is not available in API mode"), true
	}

	handled = true
	err = opts.HandleScrape() // Note this never returns
	// EXISTING_CODE

	return
}

// EXISTING_CODE
// EXISTING_CODE
