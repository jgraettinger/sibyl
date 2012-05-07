package cclparse

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

	if !baseline.Equals(compare) {
		t.Error("Invariant to argument declaration order")
	}

	compare = NewParseNode(C_goes,
		).AddLabelArgument("err", NewParseNode(C_it),
		).AddLabelArgument("l2", NewParseNode(C_here))

	if baseline.Equals(compare) {
		t.Error("Label mismatch")
	}
}

func TestBuildDirectedParse(t *testing.T) {

	// simple case 
	{
		chart := NewChart()

		C_the := chart.AddCell("the")
		C_boy := chart.AddCell("boy")

		/*D_the_boy :=*/ chart.AddLink(0, 1, 0)
		/*D_boy_the :=*/ chart.AddLink(1, 0, 0)

		expected := NewParseNode(C_the, C_boy)
		parse := chart.BuildDirectedParse()

		if !parse.Equals(expected) {
			t.Error(parse, expected)
		}
	}

	// complex case
	{
		chart := NewChart()

		C_I := chart.AddCell("I")
		C_know := chart.AddCell("know")
		C_the := chart.AddCell("the")
		C_boy := chart.AddCell("boy")
		C_sleeps := chart.AddCell("sleeps")

		D1_know_I := chart.AddLink(1, 0, 1)
		D0_know_the := chart.AddLink(1, 2, 0)
		/*D0_the_boy :=*/ chart.AddLink(2, 3, 0)
		/*D0_boy_the :=*/ chart.AddLink(3, 2, 0)
		D0_know_sleeps := chart.AddLink(1, 4, 0)
		D1_sleeps_boy := chart.AddLink(4, 3, 1)

		PN_the_boy := NewParseNode(C_the, C_boy)

		expected := NewParseNode(C_know,
			).AddLinkArgument(D1_know_I, NewParseNode(C_I),
			).AddLinkArgument(D0_know_the, PN_the_boy,
			).AddLinkArgument(D0_know_sleeps, NewParseNode(C_sleeps,
				).AddLinkArgument(D1_sleeps_boy, PN_the_boy))

		head := chart.BuildDirectedParse()

		if !head.Equals(expected) {
			t.Error(head, expected)
		}
	}
}

