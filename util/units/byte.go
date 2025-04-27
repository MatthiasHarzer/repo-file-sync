package units

import "fmt"

const (
	Byte int64 = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
)

func formatSize(size int64, unit int64, unitName string) string {
	value := float64(size) / float64(unit)
	return fmt.Sprintf("%.*f %s", 2, value, unitName)
}

func ConvertBytesToHumanReadable(size int64) string {
	switch {
	case size >= TiB:
		return formatSize(size, TiB, "TiB")
	case size >= GiB:
		return formatSize(size, GiB, "GiB")
	case size >= MiB:
		return formatSize(size, MiB, "MiB")
	case size >= KiB:
		return formatSize(size, KiB, "KiB")
	default:
		return formatSize(size, Byte, "Byte")
	}
}
