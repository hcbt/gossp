package vocoder

import (
	"github.com/hcbt/gossp"
	"github.com/hcbt/gossp/excite"
	"github.com/hcbt/gossp/f0"
	"github.com/hcbt/gossp/io"
	"github.com/hcbt/gossp/mgcep"
	"github.com/hcbt/gossp/window"
	"log"
	"testing"
)

func TestMGLSASynthesis(t *testing.T) {
	var (
		testData   []float64
		frameShift = 80
		frameLen   = 512
		alpha      = 0.41
		stage      = 12
		gamma      = -1.0 / float64(stage)
		order      = 24
		f0Seq      []float64
		ex         []float64
		mgc        [][]float64
	)

	w, err := io.ReadWav("../test_files/test16k.wav")
	if err != nil {
		log.Fatal(err)
	}
	testData = w.GetMonoData()

	// F0
	f0Seq = f0.SWIPE(testData, 16000, frameShift, 60.0, 700.0)

	// MGCep
	frames := gossp.DivideFrames(testData, frameLen, frameShift)
	mgc = make([][]float64, len(frames))
	for i, frame := range frames {
		mgc[i] = mgcep.MGCep(window.BlackmanNormalized(frame),
			order, alpha, gamma)
	}

	// adjast number of frames
	m := min(len(f0Seq), len(mgc))
	f0Seq, mgc = f0Seq[:m], mgc[:m]

	// Excitation
	g := &excite.PulseExcite{
		SampleRate: 16000,
		FrameShift: frameShift,
	}
	ex = g.Generate(f0Seq)

	// Waveform generation
	synth := NewMGLSASpeechSynthesizer(order, alpha, stage, frameShift)

	_ = synth.Synthesis(ex, mgc)
	// TODO(ryuichi) valid check
}
