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
#include "traverser.h"

namespace qblocks {

//-----------------------------------------------------------------------
bool forEveryAppearance(const CTraverserArray& traversers, const CAppearanceArray_mon& apps, void* data) {
    // If we have nothing to do, return success
    if (!apps.size() || !traversers.size())
        return true;

    // for each traverser...
    for (auto trav : traversers) {
        if (shouldQuit())
            break;
        if (!trav.traverse(apps, data))
            return false;
    }

    return true;
};

//-----------------------------------------------------------------------
bool CTraverser::traverse(const CAppearanceArray_mon& apps, void* data) {
    // Prepare the export...
    if (preFunc && !(*preFunc)(this, data))
        return false;

    blknum_t start = 0;
    if (traverserRange.first != 0) {
        for (index = 0; index < apps.size() && !shouldQuit(); index++) {
            if (apps[index].blk >= traverserRange.first) {
                break;
            }
            start = index;
        }
    }

    // For each appearance...
    for (index = start; index < apps.size() && !shouldQuit(); index++) {
        app = &apps[index];
        trans = CTransaction();  // reset
        bool passedFilter = !filterFunc || (*filterFunc)(this, data);
        if (passedFilter) {
            if (dataFunc && !(*dataFunc)(this, data))
                return (!postFunc || (*postFunc)(this, data)) && false;
            if (displayFunc && !(*displayFunc)(this, data))
                return (!postFunc || (*postFunc)(this, data)) && false;
            nProcessed += (counterFunc ? (*counterFunc)(this, data) : 1);
        }
    }

    // Cleanup the export...
    if (postFunc && !(*postFunc)(this, data))
        return false;

    return true;
}

//-----------------------------------------------------------------------
bool filterByRange(CTraverser* trav, void* data) {
    if (!bn_2_Timestamp(trav->app->blk) || shouldQuit())
        return false;
    return inRange(blknum_t(trav->app->blk), trav->traverserRange.first, trav->traverserRange.second);
}

//-------------------------------------------------------------------------
ostream& operator<<(ostream& os, const CTraverser& it) {
    os << "logging: " << it.logging << endl;
    os << "index: " << it.index << endl;
    os << "nProcessed: " << it.nProcessed << endl;
    os << "searchType: " << it.searchType << endl;
    os << "curMonitor: " << it.curMonitor << endl;
    // os << "monitorMap: " << it.monitorMap << endl;
    os << "traverserRange: " << it.traverserRange << endl;
    return os;
}

}  // namespace qblocks
