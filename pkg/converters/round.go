package converters

import "math"

func Round(amount float64) float64 {
	return math.Round(amount*100) / 100
}
