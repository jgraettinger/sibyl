package cclparse

import (
    "fmt"
)

type ParseNode struct {
    From *Cell // inclusive
    To   *Cell // inclusive

    Arguments []ParseNodeArgument
}

type ParseNodeArgument struct {
    Child *Node

    Link *CoverLink
    Label string
}

func (chart *Chart) BuildDirectedParse() (head *ParseNode) {

    covering := make([]*ParseNode, len(chart))

    collapse := false
    for _, cell := range(chart) {

        if collapse {
            invariant(ind != 0, "non-zero index")

            covering[cell.Index] = covering[cell.Index - 1]
            covering[cell.Index].To = cell
        }

        var forwardLink, backLink *Link

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


        } else {
            // invent a new ParseNode to cover this cell

            covering[cell.Index] = &ParseNode{cell, cell, []ParseNodeArgument}
        }
    }

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

func ExtractDirectedCoverLinkParse(chart Chart) *ParseNode {

    var lastCell *Cell

    var lastParseNOde

    for _, cell := range(chart) {


    }
}

func ExtractConstituentParse(chart Chart) map[*Link] uint {

    heights := make(map[*Link] uint)

    for _, cell := range(chart) {
        for _, link := range(cell.Inbound) {
            annotateLink(link, heights, 1)
        }
    }
    return heights
}

func annotateLink(link *Link, heights map[*Link] uint, curHeight uint) {
    if heights[link] >= curHeight {
        return
    }
    heights[link] = curHeight
    fmt.Println(heights)

    // recursively follow back-links up to the root
    for _, parent := range(link.From.Inbound) {
        annotateLink(parent, heights, curHeight + 1)
    }
}

