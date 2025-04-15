package utils

import "math"

func RoundOfTo2DP(amount float64) float64 {
	return math.Round(amount*100) / 100
}
