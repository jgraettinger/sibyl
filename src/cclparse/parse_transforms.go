package cclparse

func (chart Chart) BuildDirectedParse() *ParseNode {

	// step one: create parse nodes, collapsing spans of adjacent
	//  directed cycles into a shared ParseNode
	covering := make([]*ParseNode, len(chart))

	collapse := false
	for _, cell := range(chart) {

		if collapse {
			// expand previous cell's cover node to this cell
			node := covering[cell.Index - 1]
			covering[cell.Index] = node
			node.covered = append(node.covered, cell)
		} else {
			// invent a new ParseNode to cover this cell
			covering[cell.Index] = NewParseNode(cell)
		}

		var forwardLink, backLink *CoverLink

		for _, link := range(cell.Outbound) {
			if link.To.Index == cell.Index + 1 {
				forwardLink = link
			}
		}
		for _, link := range(cell.Inbound) {
			if link.From.Index == cell.Index + 1 {
				backLink = link
			}
		}

		if forwardLink != nil && backLink != nil {
			// this cell forms a cycle with the previous cell;
			//  collapse into a shared ParseNode

			invariant(forwardLink.Depth == 0 && backLink.Depth == 0,
				"Bad cycle %v <-> %v", forwardLink, backLink);

			collapse = true
		} else {
			collapse = false
		}
	}

	// step 2: build ParseNodeArgument links derived from CoverLinks;
	//  also track which ParseNodes are reachable for head detection
	heads := make(map[*ParseNode] bool)

	for _, node := range(covering) {
		heads[node] = true
	}

	for _, cell := range(chart) {

		node := covering[cell.Index]
		for _, link := range(cell.Outbound) {

			child := covering[link.To.Index]
			if node == child {
				// skip self-referential links
				continue
			}

			node.AddLinkArgument(link, child)
			delete(heads, child)
		}
	}

	invariant(len(heads) == 1, "Multiple heads detected: %v", heads)
	for head := range(heads) {
		return head
	}
	// not reached
	return nil
}

func (chart Chart) BuildDependencyParse() *ParseNode {

	head := chart.BuildDirectedParse()

	// step 1: for each node, identify it's deepest link within the tree
	nodeDepths := make(map[*ParseNode] uint)

	// recursive closure to annotate nodes with max depths
	var gatherDepths func(uint, *ParseNode)
	gatherDepths = func(depth uint, node *ParseNode) {
		if depth > nodeDepths[node] {
			nodeDepths[node] = depth
		}
		for child := range(node.arguments) {
			gatherDepths(depth + 1, child)
		}
	}
	gatherDepths(0, head)

	// step 2: remove all but the deepest links to a node
	var removeLinks func(uint, *ParseNode)
	removeLinks = func(depth uint, node *ParseNode) {
		depth += 1
		for child, _ := range(node.arguments) {
			if depth < nodeDepths[child] {
				// there's a deeper link to this child within
				//  the parse; remove this one
				delete(node.arguments, child)
			} else {
				removeLinks(depth, child)
			}
		}
	}
	removeLinks(0, head)

	return head
}

func (chart Chart) BuildConstituentParse() *ParseNode {

	head := chart.BuildDependencyParse()

	var hoist func(*ParseNode) *ParseNode
	hoist = func(node *ParseNode) *ParseNode {

		var invented *ParseNode
		newArguments := make(ParseLinks)

		for child, link := range(node.arguments) {
			child := hoist(child)

			// update parent to reference invented node
			if link.coverLink != nil && link.coverLink.Depth == 1 {

				// link from a d=1 cover-link; invent a covering node
				//  with the same head, and host argument to it 
				if invented == nil {
					invented = NewParseNode(node.covered...)
					invented.AddLabelArgument("invented", node)
				}
				invented.AddLinkArgument(link.coverLink, child)

			} else {
				newArguments[child] = link
			}
		}
		node.arguments = newArguments

		if invented != nil {
			return invented
		}
		return node
	}
	head = hoist(head)
	return head
}

