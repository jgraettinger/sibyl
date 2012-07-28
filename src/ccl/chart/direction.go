package chart

import (
    "invariant"
)

type Direction bool

const (
	LEFT_TO_RIGHT Direction = false
	RIGHT_TO_LEFT Direction = true
)

func DirectionFromPosition(position int) Direction {
	if position < 0 {
		return RIGHT_TO_LEFT
	}
	return LEFT_TO_RIGHT
}

func (dir Direction) Flip() Direction {
	return !dir
}

func (dir Direction) Side() int {
    if dir == LEFT_TO_RIGHT {
        return LEFT
    }
    return RIGHT
}

func (dir Direction) String() string {
    if dir == LEFT_TO_RIGHT {
        return "LEFT_TO_RIGHT"
    }
    return "RIGHT_TO_LEFT"
}

/*
// Returns index of first token in this direction
func (dir Direction) BeginIndex(chart *Chart) int {
	if dir == LEFT_TO_RIGHT {
		return 0
	}
	return len(chart.Cells) - 1
}

// Returns one past the last token in this direction
func (dir Direction) EndIndex(chart *Chart) int {
	if dir == LEFT_TO_RIGHT {
		return len(chart.Cells)
	}
	return -1
}

func (dir Direction) AtEnd(index int, chart *Chart) bool {
	if dir == LEFT_TO_RIGHT {
		return index == len(chart.Cells)
	}
	return index == -1
}
*/

func (dir Direction) Increment(index int) int {
	if dir == LEFT_TO_RIGHT {
		return index + 1
	}
	return index - 1
}

func (dir Direction) Decrement(index int) int {
	if dir == LEFT_TO_RIGHT {
		return index - 1
	}
	return index + 1
}

func (dir Direction) Less(index1, index2 int) bool {
	if dir == LEFT_TO_RIGHT {
		return index1 < index2
	}
	return index1 > index2
}

func (dir Direction) OutboundAdjacency(cell *Cell) *Adjacency {
	if dir == LEFT_TO_RIGHT {
		return cell.OutboundAdjacency[RIGHT]
	}
	return cell.OutboundAdjacency[LEFT]
}

func (dir Direction) SetOutboundAdjacency(cell *Cell, adjacency *Adjacency) {
    if dir == LEFT_TO_RIGHT {
        cell.OutboundAdjacency[RIGHT] = adjacency
    } else {
        cell.OutboundAdjacency[LEFT] = adjacency
    }
}

func (dir Direction) InboundAdjacencies(cell *Cell) AdjacencySet {
    if dir == LEFT_TO_RIGHT {
        return cell.InboundAdjacencies[LEFT]
    }
    return cell.InboundAdjacencies[RIGHT]
}

func (dir Direction) OutboundLinks(cell *Cell) *LinkList {
	if dir == LEFT_TO_RIGHT {
		return &cell.OutboundLinks[RIGHT]
	}
	return &cell.OutboundLinks[LEFT]
}

func (dir Direction) InboundLink(cell *Cell) *Link {
	if dir == LEFT_TO_RIGHT {
		return cell.InboundLink[LEFT]
	}
	return cell.InboundLink[RIGHT]
}

func (dir Direction) SetInboundLink(cell *Cell, link *Link) {
    if dir == LEFT_TO_RIGHT {
        invariant.IsNil(cell.InboundLink[LEFT])
        cell.InboundLink[LEFT] = link
    } else {
        invariant.IsNil(cell.InboundLink[RIGHT])
        cell.InboundLink[RIGHT] = link
    }
}

func (dir Direction) LastOutboundLinkD0(cell *Cell) *Link {
    if dir == LEFT_TO_RIGHT {
        return cell.LastOutboundLinkD0[RIGHT]
    }
    return cell.LastOutboundLinkD0[LEFT]
}

func (dir Direction) SetLastOutboundLinkD0(cell *Cell, link *Link) {
    if dir == LEFT_TO_RIGHT {
        cell.LastOutboundLinkD0[RIGHT] = link
    } else {
        cell.LastOutboundLinkD0[LEFT] = link
    }
}

func (dir Direction) LastOutboundLinkD1(cell *Cell) *Link {
    if dir == LEFT_TO_RIGHT {
        return cell.LastOutboundLinkD1[RIGHT]
    }
    return cell.LastOutboundLinkD1[LEFT]
}

func (dir Direction) SetLastOutboundLinkD1(cell *Cell, link *Link) {
    if dir == LEFT_TO_RIGHT {
        cell.LastOutboundLinkD1[RIGHT] = link
    } else {
        cell.LastOutboundLinkD1[LEFT] = link
    }
}


/*

func (dir Direction) HasFullyBlockedAfter(cell *Cell) bool {
    if dir == LEFT_TO_RIGHT {
        return cell.FullyBlockedAfter[LEFT] != nil
    }
    return cell.FullyBlockedAfter[RIGHT] != nil
}

func (dir Direction) FullyBlockedAfter(cell *Cell) int {
    if dir == LEFT_TO_RIGHT {
        return *cell.FullyBlockedAfter[LEFT]
    }
    return *cell.FullyBlockedAfter[RIGHT]
}

func (dir Direction) SetFullyBlockedAfter(cell *Cell, index int) {
    if dir == LEFT_TO_RIGHT {
        cell.FullyBlockedAfter[LEFT] = &index
    } else {
        cell.FullyBlockedAfter[RIGHT] = &index
    }
}
*/
/*
func (dir Direction) PathBegin(cell *Cell) int {
    return cell.PathBegin[dir.Side()]
}
func (dir Direction) SetPathBegin(cell *Cell, index int) {
    cell.PathBegin[dir.Side()] = index
}
func (dir Direction) PathEndD0(cell *Cell) int {
    return cell.PathEndD0[dir.Side()]
}
func (dir Direction) SetPathEndD0(cell *Cell, index int) {
    cell.PathEndD0[dir.Side()] = index
}

func (dir Direction) PathEndD1(cell *Cell) int {
    return cell.PathEndD1[dir.Side()]
}
func (dir Direction) SetPathEndD1(cell *Cell, index int) {
    cell.PathEndD1[dir.Side()] = index
}

func (dir Direction) PathEnd(cell *Cell) int {
    return dir.Largest(dir.PathEndD0(cell), dir.PathEndD1(cell))
}
*/

func (dir Direction) Largest(a, b int) int {
    if dir == LEFT_TO_RIGHT {
        if a > b {
            return a
        }
        return b
    } else if b > a {
        return a
    }
    return b
}



/*
func (dir Direction) EnumeratePaths(cell *Cell,
	callback func(*Adjacency)) {

	for _, adjacency := range *dir.Outbound(cell) {
		if adjacency.Used {
			callback(adjacency)

			if adjacency.To != nil {
				dir.EnumeratePaths(adjacency.To, callback)
			}
		}
	}
}
*/
