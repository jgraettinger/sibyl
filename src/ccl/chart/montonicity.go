package chart

import (
	"log"
)

func (adjacency *Adjacency) MontonicityRestriction() DepthRestriction {
	forward := DirectionFromPosition(adjacency.Position)

	// 3.2.1 Montonicity: New adjacencies must have depth >= previous ones.
	if link := forward.LastOutboundLinkD1(adjacency.From); link != nil {
		log.Print("Montonicity restricts D0")
		return RESTRICT_D0
	}
	return RESTRICT_NONE
}
