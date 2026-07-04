package skills

import "math"

func CalculateCBM(dim Dimensions, quantity int) float64 {
	cbm := dim.Length * dim.Width * dim.Height / 1000000 * float64(quantity)
	return math.Round(cbm*1000) / 1000
}
