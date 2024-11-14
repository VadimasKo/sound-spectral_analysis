package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/components"
	"io"
	"modulation_effect/audio"
	"modulation_effect/chart"
	"modulation_effect/filePicker"
	"os"
	"sync"
)

var outputPath = "./output"
var frameDuration = 20.0 / 1000.0 // 20ms
var energyThreshold = 0.2

func main() {
	audioFilePath := filePicker.PickWavFile()

	file, err := os.Open(audioFilePath)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	ap, err := audio.StartAudioProcessing(file, audioFilePath)
	if err != nil {
		fmt.Printf("failed to decode .wav: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		ap.ProcessAudioFile()
	}()

	processedPoints := make([][]audio.AudioInputPoint, ap.FileProperties.ChannelCount)
	for i := range processedPoints {
		processedPoints[i] = make([]audio.AudioInputPoint, 2000)
	}

	wg.Wait()
	for channeledPoints := range ap.PointsChannel {
		for channelId, points := range channeledPoints {
			for _, point := range points {
				processedPoints[channelId] = append(processedPoints[channelId], point)
			}
		}
	}

	timeChartOptions := chart.ChartOptions{
		Title:           "Time series chart",
		ChanneledPoints: processedPoints,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	segmentChannel := make([][]audio.AudioInputPoint, 1)
	segmentChannel[0] = audio.GetSignalSegment(2.0, 0.02, float64(ap.FileProperties.SampleRate), processedPoints[0])
	segmentChartOptions := chart.ChartOptions{
		Title:           "Segment time chart",
		ChanneledPoints: segmentChannel,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	windowChannel := make([][]audio.AudioInputPoint, 1)
	windowChannel[0] = audio.ApplyWindow(segmentChannel[0])
	windowChartOptions := chart.ChartOptions{
		Title:           "Segment (window) time chart",
		ChanneledPoints: windowChannel,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	freqChannel := make([][]audio.AudioInputPoint, 1)
	freqChannel[0] = audio.ComputeAmplitudeSpectrum(windowChannel[0], float64(ap.FileProperties.SampleRate))
	freqChartOptions := chart.ChartOptions{
		Title:           "Frequency chart",
		ChanneledPoints: freqChannel,
		Segments:        make([][]audio.AudioSegment, 0),
	}

	page := components.NewPage()
	page.AddCharts(
		timeChartOptions.CreateAudioLineChart(),
		segmentChartOptions.CreateAudioLineChart(),
		windowChartOptions.CreateAudioLineChart(),
		freqChartOptions.CreateFrequencyChart(),
	)

	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.Create(outputPath + "/chart.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
