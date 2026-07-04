package skills

import (
	"math"
	"strconv"
)

func PlanContainer(totalCBM float64) string {
	switch {
	case totalCBM > 0 && totalCBM < 15:
		return "LCL 拼箱"
	case totalCBM <= 28:
		return "1\u00d720GP"
	case totalCBM <= 58:
		return "1\u00d740GP"
	case totalCBM <= 68:
		return "1\u00d740HQ"
	default:
		return strconv.Itoa(int(math.Ceil(totalCBM/68))) + "\u00d740HQ"
	}
}
