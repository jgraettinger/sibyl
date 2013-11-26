package chart

import (
	"log"

	"github.com/dademurphy/sibyl/invariant"
)

func updateResolution(chart *Chart, link *Link) {
	if DirectionFromPosition(link.Position) == LEFT_TO_RIGHT {
		updateResolutionLeftToRight(chart, link)
	} else {
		updateResolutionRightToLeft(chart, link)
	}
}

func updateResolutionLeftToRight(chart *Chart, link *Link) {
	// This link must have completely resolved all violations,
	// if there were any, because:
	// - Parser incrementality means there weren't any violations
	//   prior to adding the last cell.
	// - A left-to-right link must end at the last cell.
	if chart.minimalViolation != nil {
		invariant.IsTrue(link.From == chart.minimalViolation)
		log.Printf("right-link %v resolves minimal violation %v",
			link, chart.minimalViolation)
		chart.minimalViolation = nil
	}
}

func updateResolutionRightToLeft(chart *Chart, link *Link) {
	if chart.minimalViolation != nil && 
		(*link.FurthestPath).Index <= chart.minimalViolation.Index {
		// This link resolves an earlier resolution violation,
		// as the triggering left-to-right path has been covered.
		log.Printf("left-link %v resolves minimal violation %v",
			link, chart.minimalViolation)
		chart.minimalViolation = nil
	}

	// Determine whether this link has created a new violation.
	// Look for a backwards inbound link into the furthest cell reachable
	// along this link path.
	if inbound := LEFT_TO_RIGHT.InboundLink(*link.FurthestPath);
		inbound != nil {

		if cell := inbound.From; !cell.LinkPathReaches(link.From) {
			// This is a resolution violation. inbound.From & Link.From each
			// have paths to inbound.To, but neither cell has a link path to
			// the other.

			// Is this a minimal violation? Per 3.2.6, a minimal violation has
			// the smallest possible span, and the smallest possible depth.

			// It's a little unintuitive at first, that we'll always track the
			// violation closest to the chart end. Consider chart V, W, X, Y, Z
			// with V =0> W, W =0> X, V =0> Y, and Z =0> Y. We'll first track
			// a violation at V, d=0. Suppose we add a link Z =0> X. We'll now
			// track a violation at W, d=0. What happened?
			// 
			// First, it's not that the violation from V has really been
			// forgotten, because we'll see it again due to the inbound
			// V =0> W. Intuitively, we've discovered that we need to resolve
			// W prior to resolving V: Because W's adjacency is covered
			// by the V => Y link, the only possible resolution is Z =0> W.
			// Eg, we're forcing the use of Z's adjacency to resolve the
			// <W, Z> violation which would otherwise remain if we naively
			// used the V =0> Z adjacency at this point.
			//
			// For extra fun, note that if the V => Y link was d=1, then by
			// using the Z =0> X link we're required to ultimately use a
			// Z =0> V link (eg, Z dominates). This is because the V =0> W
			// link requires that a resolving link be d=0, which isn't
			// possible from V due to montonicity.

			log.Printf("left-link %v created violation to %v", link, cell)

			if chart.minimalViolation == nil ||
				!RIGHT_TO_LEFT.Less(chart.minimalViolation.Index, cell.Index) {

				// Montonicity guarantees that we'll see d=0 violations
				// from a cell only after all d=1 violations. Ie, just
				// track the current inbound link depth.
				invariant.IsTrue(chart.minimalViolation != cell ||
					inbound.Depth <= chart.minimalViolationDepth)

				chart.minimalViolation = cell
				chart.minimalViolationDepth = inbound.Depth

				log.Printf("volation %v is minimal (depth %v)", cell, inbound.Depth)
			}
		}
	}
}

func (adjacency *Adjacency) ResolutionRestriction(
	chart *Chart) DepthRestriction {

	if chart.minimalViolation == nil {
		return RESTRICT_NONE
	}

	var restrictTo uint
	if cell := chart.CurrentCell(); adjacency.From == cell && adjacency.Position < 0 {
		invariant.IsTrue(len(cell.OutboundLinks[LEFT]) > 0)
		restrictTo = cell.OutboundLinks[LEFT].Last().Depth
	} else if adjacency.From == chart.minimalViolation && adjacency.To == cell {
		restrictTo = chart.minimalViolationDepth
	} else {
		return RESTRICT_ALL
	}

	if restrictTo == 0 {
		return RESTRICT_D1
	}
	return RESTRICT_D0
}
