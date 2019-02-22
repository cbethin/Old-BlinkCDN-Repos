package main

import (
	"fmt"

	"./blink"
	"./blink/blinkEstimators"
)

var EstimatorFunction blink.LatencyEstimator

func main() {
	EstimatorFunction := blinkEstimators.ProbabilisticEstimator
	latency := EstimatorFunction(10.0, 10.2)
	latency = EstimatorFunction(latency, 7)
	latency = EstimatorFunction(latency, 11)
	for i := 0; i < 1000; i++ {
		latency = EstimatorFunction(latency, 10.2)
	}

	fmt.Println(latency)
	fmt.Println(blink.EstimatorFunction == nil)
}
