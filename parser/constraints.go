package parser

type DepthRestriction uint

const (
	RESTRICT_NONE DepthRestriction = 0
	RESTRICT_D0   DepthRestriction = 1
	RESTRICT_D1   DepthRestriction = 2
	RESTRICT_ALL  DepthRestriction = 3
)

type Constraint interface {
	Name() string
	RestrictedDepths(*Adjacency) DepthRestriction
}

// Interface used by stateful constraints, which need to be kept
// aprised of used links to manage tracking structures.
type LinkObserver interface {
	Observe(*Link)
}

// Interface used by constraints which need to restrict the chart's
// NextCell() operation (ie, because of resolution violations).
type NextCellRestrictor interface {
	IsNextCellRestricted() bool
}

