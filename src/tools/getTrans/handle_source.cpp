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
#include "options.h"

size_t depth = 1;
CAddressBoolMap theMap;
size_t w = 25;
string_q arrow(const string_q& addr) {
    return padLeft(string_q(depth * 2, '-') + "> ", w) + string_q(depth, ' ') + addr;
}

blknum_t firstGnosis = 19747830;
blknum_t firstMainnet = 13868853;
//--------------------------------------------------------------
bool visitSource(CTransaction& trans, void* data) {
    cout << bBlack << string_q(120, '-') << endl;
    cout << bBlack << trans.blockNumber << "." << trans.transactionIndex << " " << trans.hash << endl;
    cout << bBlack << string_q(120, '-') << cOff << endl;
    COptions* opt = reinterpret_cast<COptions*>(data);

    if (theMap[opt->statementManager.accountedFor]) {
        cout << bRed << "Skipping -- seen it before: " << opt->statementManager.accountedFor << cOff << endl;
        return true;
    }
    theMap[opt->statementManager.accountedFor] = true;

    cout << arrow(trans.to) << " " << trans.blockNumber << endl;
    opt->statementManager.getStatements(trans);

    for (auto statement : trans.statements) {
        if (theMap[statement.sender]) {
            cout << bRed << "Skipping " << statement.sender << cOff << endl;
            continue;
        }
        theMap[statement.sender] = true;

        string_q ret = doCommand("chifra list " + statement.sender + " --count --no_header --chain " + getChain() +
                                 " 2>/dev/null | cut -f2");
        cout << ret << endl;
        size_t cnt = str_2_Uint(ret);
        cout << bGreen << "Found " << cnt << " records for address " << statement.sender << endl;
        cout << "Sleeping..." << endl;
        ::sleep(1);
        if (cnt > 2000) {
            cout << bTeal << "Skipping " << statement.sender << cOff << " because too big: " << cnt << endl;
            ::sleep(4);
            continue;
        }

        depth++;
        ostringstream exportCmd;
        exportCmd << "chifra export --accounting --statements --flow in --cache --chain " << getChain();
        exportCmd << " --last_block " << trans.blockNumber;
        exportCmd << " --first_block " << (getChain() == "gnosis" ? firstGnosis : firstMainnet);
        exportCmd << " " << statement.sender;  //  << " 2>/dev/null";
        cout << exportCmd.str() << endl;
        string_q contents = doCommand(exportCmd.str());
        CStringArray lines, fields;
        explode(lines, contents, '\n');

        /*
blockNumber	transactionIndex	logindex	transactionHash	timestamp	date
assetAddress	assetSymbol	decimals	spotPrice	priceSource	accountedFor	sender
recipient	begBal	amountNet	endBal	encoding	signature	reconciliationType reconciled totalIn
amountIn	internalIn	selfDestructIn	minerBaseRewardIn	minerNephewRewardIn	minerTxFeeIn
minerUncleRewardIn	prefundIn	totalOut	amountOut	internalOut	selfDestructOut gasOut
totalOutLessGas	prevAppBlk	prevBal	begBalDiff	endBalDiff	endBalCalc
        */
        size_t nLines = 0;
        for (auto line : lines) {
            if (fields.size() == 0) {
                explode(fields, line, '\t');
                continue;
            } else {
                CStringArray values;
                explode(values, line, '\t');
                if (values.size() > 13) {
                    hash_t hash = values[3];
                    address_t sender = values[13];
                    CTransaction tr;
                    getTransaction(tr, hash);
                    opt->statementManager.getTransfers(tr);
                    for (auto t : opt->statementManager.transfers) {
                    }
                    opt->statementManager.accountedFor = opt->statementManager.transfers[0].sender;
                    cout << nLines++ << ". " << values[0] << "." << values[1] << "-" << values[2] << " " << hash << " ["
                         << tr.hash << "]" << endl;
                    if (tr.hash != "") {
                        visitSource(tr, data);
                    }
                }
            }
        }
        depth--;
    }

    return true;
}

// chifra export 0x45233ff6bbc10f0eb41f430561ac97d78ddfb4cd 2>/dev/null | cut -f3,4,10
// chifra list 0x45233ff6bbc10f0eb41f430561ac97d78ddfb4cd | wc