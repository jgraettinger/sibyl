package parser

type Chart struct {
	input TokenStream

	// Processed cells in token order.
	Cells []*Cell

	// Active adjacencies whose Tail is the last cell.
	leftToRightAdjacencies []*Adjacency

	// Adjacencies to {end}. While not valid adjacencies in this state,
	// they may become a valid adjacency to the next token processed.
	endAdjacencies []*Adjacency
}

func (chart *Chart) NextCell() (*Cell, error) {
	token, err := chart.input.NextToken()
	if err != nil {
		return nil, err
	}

	var prev, next *Cell
	if len(chart.Cells) > 0 {
		prev = chart.Cells[len(chart.Cells)-1]
	}
	next = &Cell{Index: len(chart.Cells), Token: token}
	chart.Cells = append(chart.Cells, next)

	// Adjacencies currently in leftToRightAdjacencies are no longer useable.
	// Move adjacencies to {end} to instead be adjacent to next.
	chart.leftToRightAdjacencies = chart.leftToRightAdjacencies[:0]
	for _, adjacency := range chart.endAdjacencies {
		invariant(adjacency.Tail == nil)
		adjacency.Tail = next
		adjacency.appendTo(&chart.leftToRightAdjacencies)
	}
	chart.endAdjacencies = chart.endAdjacencies[:0]

	// Create a forwards adjacency from next => {end}.
	adjacency := &Adjacency{
		Head:     next,
		Tail:     nil,
		Position: 1}
	adjacency.appendTo(&chart.endAdjacencies)
	next.Right.OutboundAdjacency = adjacency

	// Create a backwards adjacency from next => prev.
	// If prev == nil, this is an adjacency to {begin}.
	adjacency = &Adjacency{
		Head:     next,
		Tail:     prev,
		Position: -1}
	next.Left.OutboundAdjacency = adjacency
	return next, nil
}

func (chart *Chart) UseAdjacency(adjacency *Adjacency, depth int) *Link {
	if adjacency.Position < 0 {
		return chart.linkRightToLeft(adjacency, depth)
	} else {
		return chart.linkLeftToRight(adjacency, depth)
	}
}

func (chart *Chart) linkLeftToRight(usedAdjacency *Adjacency, depth int) *Link {
	head, tail := usedAdjacency.Head, usedAdjacency.Tail
	link := NewLink(usedAdjacency, depth)

	// By incrementalness of the parser, a left-to-right link is
	// to the last cell, which is also the furthest link path.
	invariant(tail == lastCell(chart.Cells))
	updateBoxedPathLeftToRight(link)
	invariant((*link.BoxedFurthestPath).Index == len(chart.Cells)-1)

	// Update link-tracking structures.
	invariant(tail.Left.InboundLink == nil)
	link.appendTo(&head.Right.OutboundLinks)
	tail.Left.InboundLink = link

	// Examine other left-to-right adjacencies into the last cell. Those spanning
	// the new link are still viable, and the adjacency is moved to {end}. Those
	// not spanning the link have been invalidated by the minimality condition,
	// and are discarded.
	for _, adjacency := range chart.leftToRightAdjacencies {
		if adjacency.Head.Index < usedAdjacency.Head.Index {
			adjacency.Tail = nil
			adjacency.appendTo(&chart.endAdjacencies)
		} else {
			// Otherwise, discard; note usedAdjacency fails the condition.
			adjacency.Head.Right.OutboundAdjacency = nil
		}
	}
	chart.leftToRightAdjacencies = chart.leftToRightAdjacencies[:0]

	// Use of this adjacency creates a new one, spanning to {end}.
	createdAdjacency := &Adjacency{
		Head:     head,
		Position: usedAdjacency.Position + 1}
	createdAdjacency.appendTo(&chart.endAdjacencies)
	head.Right.OutboundAdjacency = createdAdjacency

	return link
}

func (chart *Chart) linkRightToLeft(usedAdjacency *Adjacency, depth int) *Link {
	head, tail := usedAdjacency.Head, usedAdjacency.Tail
	link := NewLink(usedAdjacency, depth)

	invariant(head == lastCell(chart.Cells))
	updateBoxedPathRightToLeft(link)

	// Use of this adjacency creates a new one, spanning to the cell prior
	// to the furthest path along this link. If tail is the first cell,
	// this becomes an adjacency to {begin}.
	createdAdjacency := &Adjacency{
		Head:     head,
		Position: usedAdjacency.Position - 1}
	if (*link.BoxedFurthestPath).Index > 0 {
		createdAdjacency.Tail = chart.Cells[(*link.BoxedFurthestPath).Index - 1]
	}
	head.Left.OutboundAdjacency = createdAdjacency

	// Update link-tracking structures.
	invariant(tail.Right.InboundLink == nil)
	link.appendTo(&head.Left.OutboundLinks)
	tail.Right.InboundLink = link
	return link
}
