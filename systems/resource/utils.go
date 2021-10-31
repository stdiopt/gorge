package resource

import "math"

func nearestPowerOf2(n int) int {
	if isPowerOf2(n) {
		return n
	}
	return int(math.Pow(2.0, math.Round(math.Log(float64(n))/math.Log(2))))
}

func isPowerOf2(n int) bool {
	return n != 0 && (n&(n-1)) == 0
}

func makeEven(n int) int {
	return n + n&1
}
