package chart

import (
	"testing"

	"launchpad.net/gocheck"
)

// Package-level test setup.
func TestPackage(t *testing.T) { gocheck.TestingT(t) }

type CellSuite struct{}
var _ = gocheck.Suite(&CellSuite{})


/*
func (s *CellSuite) TestD0LinkPathReaches(c *gocheck.C) {
	chart := NewChart()
	chart.AddTokens("X", "Y", "Z")

	X := chart.NextCell()
	Y := chart.NextCell()
	chart.Use(X.Right.OutboundAdjacency, 0)

	type tCase struct {
		side *Side
		reach int
	}

	for _, tt := range []tCase {
		{&X.Right, 1}, 
		{&X.Left, 0},
		{&Y.Right, 1},
		{&Y.Left, 0},
		{&Z.Right, 2},
		{&Z.Left, 2},
	} {
		c.Check(tt.side.FurthestD0Path(), gocheck.Equals, tt.reach)
	}

	chart.Use(Y.Right.OutboundAdjacency, 0)

	for _, tt := range []tCase {
		{&X.Right, 1}, 
		{&X.Left, 0},
		{&Y.Right, 1},
		{&Y.Left, 0},
		{&Z.Right, 2},
		{&Z.Left, 2},
	} {
		c.Check(tt.from.D0LinkPathReaches(tt.to), gocheck.Equals, tt.expect)
	}
}

func (s *CellSuite) TestD1LinkPathReaches(c *gocheck.C) {
	chart := NewChart()
	chart.AddTokens("X", "Y", "Z")

	X := chart.NextCell()
	Y := chart.NextCell()
	chart.Use(X.Right.OutboundAdjacency, 1)

	type tCase struct {
		from, to, *Cell
		expect bool
	}

	for _, tt := range []tCase {
		{X, Y, true}, 
		{X, Z, false},
		{Y, X, true},
		{Y, Z, false},
		{Z, Y, false},
		{Z, X, false},
	} {
		c.Check(tt.from.D0LinkPathReaches(tt.to), gocheck.Equals, tt.expect)
	}

	chart.Use(Y.Right.OutboundAdjacency, 1)

	for _, tt := range []tCase {
		{X, Y, true}, 
		{X, Z, true},
		{Y, X, true},
		{Y, Z, true},
		{Z, Y, true},
		{Z, X, true},
	} {
		c.Check(tt.from.D0LinkPathReaches(tt.to), gocheck.Equals, tt.expect)
	}
}
*/
