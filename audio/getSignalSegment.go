package audio

import "sort"

func GetSignalSegment(start, duration, sampleRate float64, signal []AudioInputPoint) []AudioInputPoint {
	sort.Slice(signal, func(i, j int) bool {
		return signal[i].X < signal[j].X
	})

	segment := make([]AudioInputPoint, int(duration*sampleRate))
	frameEnd := start + duration
	startIndex := len(signal) - 1

	for i, point := range signal {
		if point.X > start && point.X < frameEnd {
			if startIndex > i {
				startIndex = i
			}

			segment[i-startIndex] = point
		}
	}

	return segment
}
