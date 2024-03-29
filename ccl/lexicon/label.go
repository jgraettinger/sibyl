package lexicon

import (
	"fmt"
)

type LabelType uint8

const (
	CLASS     LabelType = 0
	ADJACENCY LabelType = 1
)

type Label struct {
	Type  LabelType
	Token string
}

func (label Label) Flip() Label {
	if label.Type == CLASS {
		return Label{ADJACENCY, label.Token}
	}
	return Label{CLASS, label.Token}
}

func (label Label) IsClass() bool {
	return label.Type == CLASS
}

func (label Label) IsAdjacency() bool {
	return label.Type == ADJACENCY
}

func (label Label) String() string {
	if label.IsClass() {
		return fmt.Sprintf("%v", label.Token)
	}
	return fmt.Sprintf(".%v", label.Token)
}
