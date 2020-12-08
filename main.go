package main 

import (
	"flag"
	"fmt"
	"github.com/hcbt/gossp"
	"github.com/hcbt/gossp/io"
	"github.com/hcbt/gossp/stft"
	"github.com/hcbt/gossp/window"
	"log"
	"math"
)

func main() {
	filename := flag.String("i", "input.wav", "Input filename")
	flag.Parse()

	w, werr := io.ReadWav(*filename)
	if werr != nil {
		log.Fatal(werr)
	}
	data := w.GetMonoData()

	s := &stft.STFT{
		FrameShift: int(float64(w.SampleRate) / 100.0), // 0.01 sec,
		FrameLen:   2048,
		Window:     window.CreateHanning(2048),
	}

	spectrogram, _ := gossp.SplitSpectrogram(s.STFT(data))
	PrintMatrixAsGnuplotFormat(spectrogram)	
}

func PrintMatrixAsGnuplotFormat(matrix [][]float64) {
	fmt.Println("#", len(matrix[0]), len(matrix)/2)
	for i, vec := range matrix {
		for j, val := range vec[:1024] {
			fmt.Println(i, j, math.Log(val))
		}
		fmt.Println("")
	}
}
