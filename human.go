package code

import "fmt"

func FormatHumanReadable(size int64) string {
	const unitSize = 1024
	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	if size < unitSize {
		return fmt.Sprintf("%dB", size)
	}

	value := float64(size)
	unit := 0
	for value >= unitSize && unit < len(units)-1 {
		value /= unitSize
		unit++
	}
	return fmt.Sprintf("%.1f%s", value, units[unit])
}
