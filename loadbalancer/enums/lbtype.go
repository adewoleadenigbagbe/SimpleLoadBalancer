package enums

type LoadBalancingAlgorithmType int

const (
	RoundRobin LoadBalancingAlgorithmType = iota + 1
	WeightedRoundRobin
	IpHash
	LeastConnection
	LeastResponseTime
	ResourceLoad
)
