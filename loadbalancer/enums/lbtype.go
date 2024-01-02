package enums

type LoadBalancingAlgorithmType int

const (
	RoundRobin LoadBalancingAlgorithmType = iota + 1
	RandomWeightedRoundRobin
	SmoothWeightedRoundRobin
	IpHash
	LeastConnection
	LeastResponseTime
	ResourceLoad
)
