package chart

/*
func UpdateBlocking(chart *Chart, usedAdjacency, newAdjacency *Adjacency) {

    forward := DirectionFromPosition(usedAdjacency.Position)
    backward := forward.Flip()

    cellFrom, cellTo := usedAdjacency.From, usedAdjacency.To

    // 3.2.1 Full blocking condition

        // Blocking requires that no adjacency may span a d=1 link from a
        // node which has also has a path to the adjacency base.

        // Step 1: Propogate blocking along link-paths.
        if usedAdjacency.UsedDepth == 1 {
            // As this is a new d=1 link, we must propogate blocking
            // implications along all link-paths rooted at cellFrom.
            for _, adjacency := range *forward.Outbound(cellFrom) {
                if adjacency.Used {
                    updateFullBlocking(chart, forward, adjacency.To, cellFrom.Index)
                }
            }
            for _, adjacency := range *backward.Outbound(cellFrom) {
                if adjacency.Used {
                    updateFullBlocking(chart, backward, adjacency.To, cellFrom.Index)
                }
            }
        } else {
            // This link represents a new link-path along which blocking
            // may need to be back-propagated (in the other direction).
            if cellFrom.PathEndD1[0] != cellFrom.Index || cellFrom.PathEndD1[1] != cellFrom.Index {
                // Existing d=1 link bounds spans over cellFrom.
                updateFullBlocking(chart, forward, cellTo, cellFrom.Index)
            } else if backward.HasFullyBlockedAfter(cellFrom) {
                // Propogate constraint from farther forward in the path.
                updateFullBlocking(chart, forward, cellTo,
                    backward.FullyBlockedAfter(cellFrom))
            }
        }

        // Step 2: If newAdjacency spans FullyBlockedAfter of cellFrom, there must
        // must be a d=1 link from a node in the span (cellFrom.Index, nextAdjacencyIndex).
        if forward.HasFullyBlockedAfter(cellFrom) &&
            forward.Less(forward.FullyBlockedAfter(cellFrom, chart.ToIndex(newAdjacency))) {

            log.Printf("Immediately fully blocking new %v (blocked after %v)",
                newAdjacency, forward.FullyBlockedAfter(cellFrom))
			newAdjacency.Blocked = [2]bool{true, true}
        }


    // 3.2.1 Partial (d=0) blocking condition

        // Step 1: If there is a used, backward, inbound adjacency into
        // cellFrom which is covered by newAdjacency, than d=0 is blocked.
		if usedInbound := backward.UsedInbound(cellFrom); usedInbound != nil &&
			forward.Less(usedInbound.From.Index, nextAdjacencyIndex) {
			log.Printf("Blocking d=0 of new adjacency because of used inbound %v",
				usedAdjacency)
			newAdjacency.Blocked[0] = true
		}

	    // Step 2: Inversely, the use of this adjacency blocks d=0 of an
		// unused, backward adjacency from cellTo spanning beyond cellFrom.
		outboundAdjacency := backward.Outbound(cellTo).Current()
		log.Printf("Considering flipped outbound %v", outboundAdjacency)
		if backward.Less(cellFrom.Index, outboundAdjacency.ToIndex()) {
			log.Printf("Blocking d=0 of %v", outboundAdjacency)
			outboundAdjacency.Blocked[0] = true
		}
}

// Follows link-paths of 'root' in the direction 'backDir', enumerating
// cells to which root has a link path. For each such cell,
// - cell.FullyBlockedAfter is updated with root's index in the forward
//   direction, capturing the constraint that no adjacency may span
//   beyond this index (due to blocking, section 3.2.1 condition 3).
// - Unused adjacencies of cell in the forward direction which span
//   root.Index are blocked (& presence of a used adjacency is an
//   invariant violation).
func updateFullBlocking(chart *Chart, backDir Direction,
	cell *Cell, blockIndex int) {

	forwardDir := backDir.Flip()
    invariant.NotNil(cell)

    invariant.NotEqual(cell.Index, blockIndex)

    // Track the minimal index after which adjacencies are fully blocked.
    // The cell which roots the d=1 path may be passed in as a convienence, to
    // prime the recursive enumeration. It shouldn't be modified itself.
    updateBlocking := cell.Index != blockIndex && (
        forwardDir.FullyBlockedAfter(cell) == nil ||
        forwardDir.Less(blockIndex, *forwardDir.FullyBlockedAfter(cell)))

    if updateBlocking {
	    log.Printf("Updating FullyBlockedAfter of cell %v, %v => %v",
            cell, forwardDir, blockIndex)
        forwardDir.SetFullyBlockedAfter(cell, blockIndex)
    }

    if cell.IndexTo != blockIndex {
        // Enumerate forward adjacencies, blocking any which span blockIndex.
	    for adjacency := range forwardDir.Outbound(cell) {
		    if forwardDir.Less(blockIndex, forwardDir.IndexTo(adjacency)) {
			    invariant.IsFalse(adjacency.Used)
                invariant.IsTrue(updateBlocking)

    			if !outAdjacency.IsBlocked() {
	    			log.Printf("Fully blocking %v (due to back-linked d=1 from %v)",
		    			adjacency, *blockIndex)
			    	adjacency.Blocked = [2]bool{true, true}
    			}
		    }
    	}
    }


    // TODO: if updateBlocking == false, I think we can stop here;
    // we shouldn't need to recurse. Note updateBlocking will be
    // true of cell.Index == blockIndex

    // Depth-first recursive call
	for backAdjacency := range backDir.Outbound(cell) {
		if !backAdjacency.Used {
			continue
		}
        chart.updateFullBlocking(backDir, backAdjacency.To, blockIndex)
	}

    if cell.Index == blockIndex {
        return
    }

    if !updateBlocking {
        return
    }

}

*/
