package blinkEstimators

import "math"

func mean(data []float64) float64 {
	sum := float64(0)
	for _, value := range data {
		sum += value
	}

	return sum / float64(len(data))
}

func sampleVariance(data []float64) float64 {
	avg := mean(data)

	sumOfSquares := float64(0)
	for _, val := range data {
		sumOfSquares += (float64(val) - avg) * (float64(val) - avg)
	}

	return (sumOfSquares / float64(len(data)-1))
}

func standardDeviationSample(data []float64) float64 {
	return math.Pow(sampleVariance(data), 0.5)
}
