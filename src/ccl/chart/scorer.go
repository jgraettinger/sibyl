package chart

type Scorer interface {
	Score(adjacency *Adjacency) (score float64, depth uint)
}
