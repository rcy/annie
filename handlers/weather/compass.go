package weather

import "math"

func compass16(deg uint) string {
	return []string{"n", "nne", "ne", "ene", "e", "ese", "se", "sse", "s", "ssw", "sw", "wsw", "w", "wnw", "nw", "nnw"}[int(math.Round(float64(deg)/22.5))%16]
}
