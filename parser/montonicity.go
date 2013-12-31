package parser

type Montonicity bool

func (_ Montonicity) Name() string {
	return "Montonicity"
}
// Under montonicity, new adjacencies must have depth >= previous ones.
func (_ Montonicity) RestrictedDepths(a *Adjacency) DepthRestriction {
	if ll := lastLink(a.HeadSide().OutboundLinks); ll != nil && ll.Depth > 0 {
		return RESTRICT_D0
	}
	return RESTRICT_NONE
}
