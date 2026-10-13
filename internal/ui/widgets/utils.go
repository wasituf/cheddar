package ui

import "math"

func EaseInOut(p, n float64) float64 {
	if p < 0.5 {
		return 0.5 * math.Pow(2*p, n)
	}
	return 1 - 0.5*math.Pow(2-2*p, n)
}

func EaseInOutExpo(p float64) float64 {
	if p == 0 || p == 1 {
		return p
	}
	if p < 0.5 {
		return math.Pow(2, 20*p-10) / 2
	}
	return (2 - math.Pow(2, -20*p+10)) / 2
}
