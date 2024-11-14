package audio

import (
	"gonum.org/v1/gonum/dsp/fourier" // Import Fourier transforms from gonum
	"math/cmplx"
)

var threshold = 0.0

func ComputeAmplitudeSpectrum(points []AudioInputPoint, sampleRate float64) []AudioInputPoint {
	signal := make([]float64, len(points))
	for i, point := range points {
		signal[i] = point.Y
	}

	N := len(signal)

	// Step 2: Compute the FFT
	fft := fourier.NewFFT(N)
	freqComplex := fft.Coefficients(nil, signal)

	// Step 3: Calculate amplitude spectrum
	amplitude := make([]float64, N/2+1)
	for i := 0; i < len(amplitude); i++ {
		amplitude[i] = cmplx.Abs(freqComplex[i])
	}

	// Step 4: Apply scaling to amplitude spectrum (except DC and Nyquist frequencies)
	for i := 1; i < len(amplitude)-1; i++ {
		amplitude[i] *= 2
	}

	// Step 5: Generate frequency axis (0 to sampleRate/2)
	frequency := make([]float64, len(amplitude))
	for i := 0; i < len(amplitude); i++ {
		frequency[i] = float64(i) * (sampleRate / float64(N)) // Frequency axis: [0 ... Fs/2]
	}

	frequencyPoints := make([]AudioInputPoint, len(frequency))
	for i := range frequency {
		if threshold < amplitude[i] {
			frequencyPoints[i].X = frequency[i]
			frequencyPoints[i].Y = amplitude[i]
		}
	}

	return frequencyPoints
}
