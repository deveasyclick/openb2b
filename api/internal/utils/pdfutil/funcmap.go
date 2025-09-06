package pdfutil

import (
	"math"
	"strconv"
	"text/template"
)

// Add your custom template functions here
var funcMap = template.FuncMap{
	"mul": func(a, b any) float64 {
		return toFloat64(a) * toFloat64(b)
	},
	"round": func(v float64) float64 {
		return math.Round(v*100) / 100
	},
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}
