package cclparse

import (
    "fmt"
)

type Cell struct {
    Index uint
    Token string

    Inbound  []*CoverLink
    Outbound []*CoverLink
}

type CoverLink struct {
    From *Cell // inclusive
    To   *Cell // inclusive
    Depth uint
}

type Chart []*Cell

func NewChart() (chart *Chart) {
    chart = &Chart{}
    return
}

func (chart *Chart) AddCell(token string) (cell *Cell) {

    cell = &Cell{(uint)(len(*chart)), token, []*CoverLink{}, []*CoverLink{}}
    *chart = append(*chart, cell)
    return cell
}

func (c Chart) AddLink(from, to, depth uint) (link *CoverLink) {
    link = &CoverLink{c[from], c[to], depth}

    c[from].Outbound = append(c[from].Outbound, link)
    c[to].Inbound = append(c[to].Inbound, link)
    return
}

func (c *Cell) String() string {
    return fmt.Sprintf("Cell<%d, %s, %v, %v>",
        c.Index, c.Token, c.Inbound, c.Outbound)
}

func (l *CoverLink) String() string {
    return fmt.Sprintf("CoverLink<%s (%d), %s (%d), %d>",
        l.From.Token, l.From.Index, l.To.Token, l.To.Index, l.Depth)
}

