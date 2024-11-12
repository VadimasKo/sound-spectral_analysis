package audio

import (
	"math"
	"sort"
)

func ApplyModulationFilter(sound [][]AudioInputPoint, props AudioFileProperties, modFrequency float64) [][]AudioInputPoint {
	channels := make([][]AudioInputPoint, len(sound))

	for i, channel := range sound {
		channels[i] = modulateChannel(channel, modFrequency)
	}

	return channels
}

func modulateChannel(channel []AudioInputPoint, modFrequency float64) []AudioInputPoint {
	sampleLength := len(channel)

	sort.Slice(channel, func(i, j int) bool {
		return channel[i].X < channel[j].X
	})

	signalArray := make([]AudioInputPoint, sampleLength)
	for i, point := range channel {
		t := point.X

		// Create the modulating waveform (sine wave)
		modulation := 0.5 * (1 + math.Sin(2*math.Pi*modFrequency*t))

		signalArray[i] = AudioInputPoint{Y: modulation * point.Y, X: point.X}
	}
	return signalArray
}
