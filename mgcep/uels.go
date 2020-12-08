package mgcep

import (
	"github.com/hcbt/gossp/3rdparty/sptk"
)

// TODO(ryuichi) replace with pure Go.
func UELS(audioBuffer []float64, order int) []float64 {
	return sptk.UELSWithDefaultParameters(audioBuffer, order)
}
