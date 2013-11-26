package cclparse

/*
type Bracket struct {
    From *Cell // inclusive
    To   *Cell // inclusive
}

type Bracketing []Bracket

func (b *Bracket) Covers(cell *Cell) bool {
    return b.From.Index <= cell.Index && b.To.Index >= cell.Index
}

func (b *Bracketing) extend(covered *Cell, extend *Cell) {
    if extend.Index < covered.Index {
        for ind := range(b) {
            // extend covering brackets leftwards as needed
            if b[ind].Covers(covered) && !b[ind].Covers(extend) {
                b[ind].From = extend
            }
        }
    } else {
        for ind := range(b) {
            // extend covering brackets rightward as needed
            if b[ind].Covers(covered) && !b[ind].Covers(extend) {
                b[ind].To = extend
            }
        }
    }
}

func (b *Bracketing) IncrementalUpdate(cell *Cell) {

    for link := range(cell.Inbound) {
        if link.Depth == 0 {
            b.extend(cell.Index, link.From.Index)
        } else if link.Depth == 1 {

        }
    }

}

func max(a, b uint32) int { if a > b { return a }; return b }
func min(a, b uint32) int { if a < b { return a }; return b }
*/
