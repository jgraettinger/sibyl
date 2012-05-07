package cclparse

import (
	"fmt"
	"sort"
	"strings"
)

type ParseNode struct {
	covered []*Cell
	arguments []ParseNodeArgument
}

type ParseNodeArgument struct {
	label string
	child *ParseNode
}

func NewParseNode(covered ...*Cell) (node *ParseNode) {
	{
		var last *Cell
		for _, cell := range(covered) {
			invariant(last == nil || last.Index + 1 == cell.Index,
				"covered cells aren't contiguous: %v", covered)
		}
	}
	node = &ParseNode{covered, []ParseNodeArgument{}}
	return
}

func (node *ParseNode) AddLinkArgument(
	link *CoverLink, child *ParseNode) *ParseNode {

	node.AddLabelArgument(fmt.Sprintf("d=%d", link.Depth), child)
	return node
}

func (node *ParseNode) AddLabelArgument(
	label string, child *ParseNode) *ParseNode {

	// small number of arguments; binary search may be overkill
	ind := sort.Search(len(node.arguments),	func (i int) bool {
			return child.covered[0].Index <
				node.arguments[i].child.covered[0].Index
		})

	node.arguments = append(node.arguments[:ind],
		append([]ParseNodeArgument{ParseNodeArgument{label, child}},
			node.arguments[ind:]...)...)

	return node
}

func (node *ParseNode) Head() string {
	var tokens []string
	for _, cell := range(node.covered) {
		tokens = append(tokens, cell.Token)
	}
	return strings.Join(tokens, " ")
}

func (node *ParseNode) Equals(other *ParseNode) bool {
	if other == nil {
		return false
	} else if len(node.covered) != len(other.covered) {
		return false
	} else if len(node.arguments) != len(other.arguments) {
		return false
	}

	for ind := range(node.covered) {
		if node.covered[ind] != other.covered[ind] {
			return false
		}
	}
	for ind := range(node.arguments) {
		if node.arguments[ind].label != other.arguments[ind].label {
			return false
		}
		if !node.arguments[ind].child.Equals(
			other.arguments[ind].child) {
			return false
		}
	}
	return true
}

func (node *ParseNode) AsText(indent int, label string) []string {

	parts := []string{
		strings.Repeat(" ", indent),
		fmt.Sprintf("<%v>(\"%v\"", label, node.Head())}

	for _, arg := range(node.arguments) {

		parts = append(parts, ",\n")
		parts = append(parts, arg.child.AsText(indent + 1, arg.label)...)
	}
	parts = append(parts, ")")
	return parts
}

func (node *ParseNode) String() string {
	return strings.Join(node.AsText(0, ""), "")
}


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

func invariant(check bool, a ...interface{}) {
	if check {
		return
	}
	if len(a) != 0 {
		errf := a[0].(string)

		panic(fmt.Sprintf(errf, a[1:]...))
	} else {
		panic("Invariant check failed")
	}
}

