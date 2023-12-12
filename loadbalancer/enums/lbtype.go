package enums

type LoadBalancingAlgorithmType int

const (
	RoundRobin LoadBalancingAlgorithmType = iota + 1
	StickyRoundRobin
	WeightedRoundRobin
	Hash
	LeastConnection
	LeastResponseTime
)
