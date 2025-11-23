package code

import "fmt"

func FormatHumanReadable(size int64) string {
	const base = 1024
	if size < base {
		return fmt.Sprintf("%dB", size)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	var value float64 = float64(size)
	var i int
	for i = 0; i < len(units); i++ {
		value = value / float64(base)
		if value < float64(base) {
			if value == float64(int64(value)) {
				return fmt.Sprintf("%.0f%s", value, units[i])
			}
			return fmt.Sprintf("%.1f%s", value, units[i])
		}
	}
	return fmt.Sprintf("%.1f%s", value, units[len(units)-1])
}