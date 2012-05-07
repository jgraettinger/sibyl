package cclparse

import (
	"fmt"
	"testing"
)

func buildFixture() (chart *Chart, C map[string]*Cell, L map[string]*CoverLink) {

	chart = NewChart()
	C = make(map[string]*Cell)
	L = make(map[string]*CoverLink)

	for _, token := range([]string{"I", "know", "the", "boy", "sleeps"}) {
		C[token] = chart.AddCell(token)
	}

	L["D1_know_I"] = chart.AddLink(1, 0, 1)
	L["D0_know_the"] = chart.AddLink(1, 2, 0)
	L["D0_the_boy"] = chart.AddLink(2, 3, 0)
	L["D0_boy_the"] = chart.AddLink(3, 2, 0)
	L["D0_know_sleeps"] = chart.AddLink(1, 4, 0)
	L["D1_sleeps_boy"] = chart.AddLink(4, 3, 1)

	return
}


func TestBuildDirectedParse(t *testing.T) {

	// single-bracket NP fragment case
	{
		chart := NewChart()
		C_a := chart.AddCell("a")
		C_purple := chart.AddCell("purple")
		C_zombie := chart.AddCell("zombie")

		chart.AddLink(0, 1, 0)
		chart.AddLink(1, 0, 0)
		chart.AddLink(1, 2, 0)
		chart.AddLink(2, 1, 0)

		expected := NewParseNode(C_a, C_purple, C_zombie)
		parse := chart.BuildDirectedParse()

		if parse.String() != expected.String() {
			t.Error(parse, expected)
		}
	}

	{
		chart, C, L := buildFixture()

		PN_the_boy := NewParseNode(C["the"], C["boy"])

		expected := NewParseNode(C["know"],
			).AddLinkArgument(L["D1_know_I"], NewParseNode(C["I"]),
			).AddLinkArgument(L["D0_know_the"], PN_the_boy,
			).AddLinkArgument(L["D0_know_sleeps"], NewParseNode(C["sleeps"],
				).AddLinkArgument(L["D1_sleeps_boy"], PN_the_boy))

		head := chart.BuildDirectedParse()

		if head.String() != expected.String() {
			t.Error(head, expected)
		}
	}
}

func TestBuildDependencyParse(t *testing.T) {
	{
		chart, C, L := buildFixture()

		expected := NewParseNode(C["know"],
			).AddLinkArgument(L["D1_know_I"], NewParseNode(C["I"]),
			).AddLinkArgument(L["D0_know_sleeps"], NewParseNode(C["sleeps"],
				).AddLinkArgument(L["D1_sleeps_boy"], NewParseNode(C["the"], C["boy"])))

		head := chart.BuildDependencyParse()
		t.Log(head)

		if head.String() != expected.String() {
			t.Error(head, expected)
		}
	}
}

func TestBuildConstituentParse(t *testing.T) {

	chart, C, L := buildFixture()

	expected := NewParseNode(C["know"],
		).AddLinkArgument(L["D1_know_I"], NewParseNode(C["I"]),
		).AddLabelArgument("invented", NewParseNode(C["know"],
			).AddLinkArgument(L["D0_know_sleeps"], NewParseNode(C["sleeps"],
				).AddLabelArgument("invented", NewParseNode(C["sleeps"]),
				).AddLinkArgument(L["D1_sleeps_boy"], NewParseNode(C["the"], C["boy"]))))

	head := chart.BuildConstituentParse()

	fmt.Println(expected.AsGraphviz())
	fmt.Println(head.AsGraphviz())

	if head.String() != expected.String() {
		t.Error(head, expected)
	}
}
