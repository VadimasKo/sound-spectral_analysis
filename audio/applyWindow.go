package audio

import "math"

// Apply window to isolate segment from adjesant signals
// Using hanning => smoothly tapper towards 0
func ApplyWindow(signal []AudioInputPoint) []AudioInputPoint {
	N := len(signal)
	windowedSignal := make([]AudioInputPoint, N)

	for i := range windowedSignal {
		window := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(N-1)))
		windowedSignal[i].Y = signal[i].Y * window
		windowedSignal[i].X = signal[i].X
	}

	return windowedSignal
}
