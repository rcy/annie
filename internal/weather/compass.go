package weather

import "math"

func compass16(deg uint) string {
	return []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}[int(math.Round(float64(deg)/22.5))%16]
}
