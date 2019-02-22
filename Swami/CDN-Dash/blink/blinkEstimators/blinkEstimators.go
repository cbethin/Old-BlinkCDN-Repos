package blinkEstimators

// WeightedAverageEstimator : creates a weighted average of the old latency value and the new
func WeightedAverageEstimator(latency float64, newLatency float64) float64 {
	return 0.9*latency + 0.1*newLatency
}

var latencies []float64

// ProbabilisticEstimator : estimates the latency by finding the value which all latency values should be less than with 98% confidence
func ProbabilisticEstimator(latency float64, newLatency float64) float64 {
	latencies = append(latencies, newLatency)

	average := mean(latencies)
	stDev := standardDeviationSample(latencies)

	// ZScores: 2.2414 (98%), 1.95996 (95%), 1.64485 (90%)
	latencyAt98 := (1.64485 * stDev) + average
	return latencyAt98
}
