package utils

import "time"

// CalcExpiresIn returns the duration as an integer millisecond count.
func CalcExpiresIn(t time.Time) int64 {
	d := time.Until(t)
	return d.Milliseconds()
}
