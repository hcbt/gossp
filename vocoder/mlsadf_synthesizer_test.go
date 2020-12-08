package vocoder

import (
	"github.com/hcbt/gossp"
	"github.com/hcbt/gossp/excite"
	"github.com/hcbt/gossp/f0"
	"github.com/hcbt/gossp/io"
	"github.com/hcbt/gossp/mgcep"
	"github.com/hcbt/gossp/window"
	"log"
	"math"
	"testing"
)

func TestMLSASynthesis(t *testing.T) {
	var (
		testData   []float64
		frameShift = 80
		frameLen   = 512
		alpha      = 0.41
		order      = 24
		pd         = 5
		f0Seq      []float64
		ex         []float64
		mc         [][]float64
	)

	w, err := io.ReadWav("../test_files/test16k.wav")
	if err != nil {
		log.Fatal(err)
	}
	testData = w.GetMonoData()

	// F0
	f0Seq = f0.SWIPE(testData, 16000, frameShift, 60.0, 700.0)

	// MCep
	frames := gossp.DivideFrames(testData, frameLen, frameShift)
	mc = make([][]float64, len(frames))
	for i, frame := range frames {
		mc[i] = mgcep.MCep(window.BlackmanNormalized(frame), order, alpha)
	}

	// adjast number of frames
	m := min(len(f0Seq), len(mc))
	f0Seq, mc = f0Seq[:m], mc[:m]

	// Excitation
	g := &excite.PulseExcite{
		SampleRate: 16000,
		FrameShift: frameShift,
	}
	ex = g.Generate(f0Seq)

	// Waveform generation
	synth := NewMLSASpeechSynthesizer(order, alpha, pd, frameShift)

	_ = synth.Synthesis(ex, mc)
	// TODO(ryuichi) valid check
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func createSin(freq float64, sampleRate, length int) []float64 {
	sin := make([]float64, length)
	for i := range sin {
		sin[i] = math.Sin(2.0 * math.Pi * freq * float64(i) / float64(sampleRate))
	}
	return sin
}
