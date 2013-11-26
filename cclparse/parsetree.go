package cclparse
/*
import (
	"fmt"
	"sort"
	"strings"
	"invariant"
)

type ParseLink struct {
	label string
	coverLink *CoverLink
}

type ParseLinks map[*ParseNode] ParseLink

type ParseLinksEntry struct {
	child *ParseNode
	link ParseLink
}

type ParseLinksEntries []ParseLinksEntry

type ParseNode struct {
	covered []*Cell
	arguments ParseLinks
}

func (l ParseLinksEntries) Len() int { return len(l) }
func (l ParseLinksEntries) Swap(i, j int) {	l[i], l[j] = l[j], l[i] }

func (l ParseLinksEntries) Less(i, j int) bool {
	// order on head's beginning token index
	return l[i].child.covered[0].Index < l[j].child.covered[0].Index
}

func NewParseNode(covered ...*Cell) (node *ParseNode) {
	var last *Cell
	for _, cell := range(covered) {
		invariant.IsTrue(last == nil || last.Index + 1 == cell.Index,
			"covered cells aren't contiguous: %v", covered)
	}
	node = &ParseNode{covered, make(ParseLinks)}
	return
}

func (node *ParseNode) AddLinkArgument(
	link *CoverLink, child *ParseNode) *ParseNode {

	node.arguments[child] = ParseLink{
		fmt.Sprintf("d=%d", link.Depth), link}
	return node
}

func (node *ParseNode) AddLabelArgument(
	label string, child *ParseNode) *ParseNode {

	node.arguments[child] = ParseLink{label, nil}
	return node
}

func (node *ParseNode) Head() string {
	var tokens []string
	for _, cell := range(node.covered) {
		tokens = append(tokens, cell.Token)
	}
	return strings.Join(tokens, " ")
}

func (node *ParseNode) OrderedArguments() ParseLinksEntries {

	entries := make(ParseLinksEntries, 0, len(node.arguments))

	for child, link := range(node.arguments) {
		entries = append(entries, ParseLinksEntry{child, link})
	}
	sort.Sort(entries)
	return entries
}

func (node *ParseNode) AsText(indent int, label string) []string {

	parts := []string{
		strings.Repeat(" ", indent),
		fmt.Sprintf("<%v>(\"%v\"", label, node.Head())}

	for _, entry := range(node.OrderedArguments()) {

		parts = append(parts, ",\n")
		parts = append(parts, entry.child.AsText(
			indent + 1, entry.link.label)...)
	}
	parts = append(parts, ")")
	return parts
}

func (head *ParseNode) AsGraphviz() string {

	parts := []string{"digraph {"}

	// interns parse nodes to a unique label for graphviz output
	gvIds := make(map[*ParseNode]string)
	gvId := func(n *ParseNode) string {
		if l, ok := gvIds[n]; ok {
			return l
		}
		gvIds[n] = fmt.Sprintf("%v%d%d", strings.Replace(n.Head(), " ", "_", -1), n.covered[0].Index, len(gvIds),
			)
		return gvIds[n]
	}

	// recursively emits labeled nodes & edges in graphviz output
	var emit func (*ParseNode)
	emit = func(n *ParseNode) {

		parts = append(parts, fmt.Sprintf("  %v [label=\"%v\"];",
			gvId(n), n.Head()))

		for child, link := range(n.arguments) {
			parts = append(parts, fmt.Sprintf("  %v -> %v [label=\"%v\"];",
				gvId(n), gvId(child), link.label))

			emit(child)
		}
	}
	emit(head)

	parts = append(parts, "}")
	return strings.Join(parts, "\n")
}

func (node *ParseNode) String() string {
	return strings.Join(node.AsText(0, ""), "")
}
*/
