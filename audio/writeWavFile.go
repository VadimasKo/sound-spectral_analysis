package audio

import (
	"fmt"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"os"
)

func WriteWavFile(filename string, props AudioFileProperties, sound [][]AudioInputPoint) error {
	// Create a new file to write the WAV data
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Create a new WAV encoder
	encoder := wav.NewEncoder(
		outFile,
		int(props.SampleRate),
		int(props.Depth),
		props.ChannelCount,
		1,
	)

	sampleCount := len(sound[0])
	audioBuffer := &audio.IntBuffer{
		Data:           make([]int, sampleCount*props.ChannelCount),
		SourceBitDepth: int(props.Depth),
		Format: &audio.Format{
			NumChannels: props.ChannelCount,
			SampleRate:  int(props.SampleRate),
		},
	}

	maxAmplitude := getMaxAmplitude(props.Depth)
	for sampleIndex := 0; sampleIndex < sampleCount; sampleIndex++ {
		for channelIndex := 0; channelIndex < props.ChannelCount; channelIndex++ {
			signal := int(sound[channelIndex][sampleIndex].Y * maxAmplitude) // Convert from float [-1, 1] to int32
			if sound[channelIndex][sampleIndex].Y > 0.95 {
				fmt.Printf("Signal channel %d, sample %d, original %f, converted %d \n", channelIndex, sampleIndex, sound[channelIndex][sampleIndex].Y, signal)
			}
			if sound[channelIndex][sampleIndex].Y < -0.95 {
				fmt.Printf("Signal channel %d, sample %d, original %f, converted %d \n", channelIndex, sampleIndex, sound[channelIndex][sampleIndex].Y, signal)
			}
			audioBuffer.Data[sampleIndex*props.ChannelCount+channelIndex] = signal
		}
	}

	if err := encoder.Write(audioBuffer); err != nil {
		return err
	}

	if err := encoder.Close(); err != nil {
		return err
	}

	return nil
}
