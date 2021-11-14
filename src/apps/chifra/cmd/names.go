/*-------------------------------------------------------------------------------------------
 * qblocks - fast, easily-accessible, fully-decentralized data from blockchains
 * copyright (c) 2016, 2021 TrueBlocks, LLC (http://trueblocks.io)
 *
 * This program is free software: you may redistribute it and/or modify it under the terms
 * of the GNU General Public License as published by the Free Software Foundation, either
 * version 3 of the License, or (at your option) any later version. This program is
 * distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even
 * the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details. You should have received a copy of the GNU General
 * Public License along with this program. If not, see http://www.gnu.org/licenses/.
 *-------------------------------------------------------------------------------------------*/
/*
 * The file was auto generated with makeClass --gocmds. DO NOT EDIT.
 */
package cmd

import (
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/cmd/globals"
	namesPkg "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/internal/names"
	"github.com/spf13/cobra"
)

// namesCmd represents the names command
var namesCmd = &cobra.Command{
	Use:   usageNames,
	Short: shortNames,
	Long:  longNames,
	RunE:  namesPkg.Run,
	Args:  namesPkg.Validate,
}

var usageNames = `names [flags] <term> [term...]

Arguments:
  terms - a space separated list of one or more search terms (required)`

var shortNames = "query addresses or names of well known accounts"

var longNames = `Purpose:
  Query addresses or names of well known accounts.`

var notesNames = `
Notes:
  - The tool will accept up to three terms, each of which must match against any field in the database.
  - The --match_case option enables case sensitive matching.`

func init() {
	namesCmd.Flags().SortFlags = false
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Expand, "expand", "e", false, "expand search to include all fields (search name, address, and symbol otherwise)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.MatchCase, "match_case", "m", false, "do case-sensitive search")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.All, "all", "l", false, "include all accounts in the search")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Custom, "custom", "c", false, "include your custom named accounts")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Prefund, "prefund", "p", false, "include prefund accounts")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Named, "named", "n", false, "include well know token and airdrop addresses in the search")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Addr, "addr", "a", false, "display only addresses in the results (useful for scripting)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Collections, "collections", "s", false, "display collections data")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Tags, "tags", "g", false, "export the list of tags and subtags only")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.ToCustom, "to_custom", "u", false, "for editCmd only, is the edited name a custom name or not (hidden)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Clean, "clean", "C", false, "clean the data (addrs to lower case, sort by addr) (hidden)")
	namesCmd.Flags().StringVarP(&namesPkg.Options.Autoname, "autoname", "A", "", "an address assumed to be a token, added automatically to names database if true (hidden)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Create, "create", "", false, "create a new name record (hidden)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Delete, "delete", "", false, "delete a name, but do not remove it (hidden)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Update, "update", "", false, "edit an existing name (hidden)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Remove, "remove", "", false, "remove a previously deleted name (hidden)")
	namesCmd.Flags().BoolVarP(&namesPkg.Options.Undelete, "undelete", "", false, "undelete a previously deleted name (hidden)")
	if os.Getenv("TEST_MODE") != "true" {
		namesCmd.Flags().MarkHidden("to_custom")
		namesCmd.Flags().MarkHidden("clean")
		namesCmd.Flags().MarkHidden("autoname")
		namesCmd.Flags().MarkHidden("create")
		namesCmd.Flags().MarkHidden("delete")
		namesCmd.Flags().MarkHidden("update")
		namesCmd.Flags().MarkHidden("remove")
		namesCmd.Flags().MarkHidden("undelete")
	}
	namesCmd.Flags().SortFlags = false
	globals.GlobalOptions(namesCmd, &namesPkg.Options.Globals)

	namesCmd.SetUsageTemplate(UsageWithNotes(notesNames))
	namesCmd.SetOut(os.Stderr)
	if namesPkg.Options.Globals.ApiMode {
		namesCmd.SetErr(os.Stdout)
	}
	chifraCmd.AddCommand(namesCmd)
}
