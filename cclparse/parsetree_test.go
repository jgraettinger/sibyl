package cclparse
/*
import (
	"testing"
)

func TestAddArgument(t *testing.T) {
	chart := NewChart()

	C_here := chart.AddCell("here")
	C_it := chart.AddCell("it")
	C_goes := chart.AddCell("goes")

	baseline := NewParseNode(C_goes,
		).AddLabelArgument("l1", NewParseNode(C_it),
		).AddLabelArgument("l2", NewParseNode(C_here))

	compare := NewParseNode(C_goes,
		).AddLabelArgument("l2", NewParseNode(C_here),
		).AddLabelArgument("l1", NewParseNode(C_it))

	if baseline.String() != compare.String() {
		t.Error("equality invariant to argument declaration order")
	}

	compare = NewParseNode(C_goes,
		).AddLabelArgument("err", NewParseNode(C_it),
		).AddLabelArgument("l2", NewParseNode(C_here))

	if baseline.String() == compare.String() {
		t.Error("Label mismatch")
	}
}
*/
