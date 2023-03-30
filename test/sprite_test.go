package sprite

import (
	"github.com/Paxx-RnD/go-ffmpeg/configuration"
	"github.com/Paxx-RnD/go-ffmpeg/ffprobe"
	"github.com/stretchr/testify/assert"
	"os"
	"sprite-preview/sprite"
	"sprite-preview/types"
	"testing"
)

var (
	inputPath      = "input.mp4"
	outputPath     = "output.png"
	vttPath        = "output.vtt"
	frameFrequency = 1
	frameWidth     = 640
	frameHeight    = 480
	numRows        = 10
	numCols        = 10
)

func Configuration() *configuration.Configuration {
	return &configuration.Configuration{
		FfprobePath: "ffprobe",
	}
}

func TestNewService(t *testing.T) {
	flags := types.Flags{}
	flags.Set()
	ffprobe := ffprobe.NewFfProbe(Configuration(), nil)
	s := sprite.NewService(flags, ffprobe)
	assert.NotNil(t, s)
}

func TestGenerateFrames(t *testing.T) {
	flags := types.Flags{Input: &inputPath}
	flags.Set()
	flags.Input = &inputPath
	ffprobe := ffprobe.NewFfProbe(Configuration(), nil)
	s := sprite.NewService(flags, ffprobe)

	frames := s.GenerateFrames()

	assert.NotEmpty(t, frames)

	for _, frame := range frames {
		_, err := os.Stat(frame)
		assert.False(t, os.IsNotExist(err))
		os.Remove(frame)
	}
}

func TestMontage(t *testing.T) {
	flags := types.Flags{Input: &inputPath, Output: &outputPath}
	flags.Set()
	flags.Input = &inputPath
	flags.Output = &outputPath
	ffprobe := ffprobe.NewFfProbe(Configuration(), nil)
	s := sprite.NewService(flags, ffprobe)

	frames := s.GenerateFrames()

	s.Montage(frames)

	_, err := os.Stat(outputPath)
	assert.False(t, os.IsNotExist(err))

	os.Remove(outputPath)
}

func TestGenerateVtt(t *testing.T) {
	flags := types.Flags{
		Input:     &inputPath,
		Output:    &outputPath,
		Vtt:       &vttPath,
		Frequency: &frameFrequency,
		Width:     &frameWidth,
		Height:    &frameHeight,
		Rows:      &numRows,
		Columns:   &numCols,
	}
	ffprobe := ffprobe.NewFfProbe(Configuration(), nil)
	s := sprite.NewService(flags, ffprobe)

	frames := s.GenerateFrames()

	s.Montage(frames)

	s.GenerateVtt(frames)

	_, err := os.Stat(vttPath)
	assert.False(t, os.IsNotExist(err))

	os.Remove(outputPath)
	os.Remove(vttPath)
}
