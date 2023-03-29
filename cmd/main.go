package main

import (
	"flag"
	"fmt"
	"github.com/Paxx-RnD/go-ffmpeg/configuration"
	"github.com/Paxx-RnD/go-ffmpeg/ffprobe"
	"github.com/Paxx-RnD/go-helper/helpers/random_helper"
	"github.com/Paxx-RnD/go-helper/helpers/time_helper"
	"os"
	"os/exec"
	"peak-tracer/types"
	"strconv"
	"time"
)

func main() {
	flags := types.Flags{}
	err := flags.Set()
	if err != nil {
		panic(err)
	}

	if flags.Help {
		flag.Usage()
		return
	}

	frames := GenerateFrames(*flags.Input, flags.Frequency, flags.Width, flags.Height)
	Montage(frames, *flags.Output)

	if !flags.GenerateVtt {
		return
	}

	GenerateVtt(frames, flags.Width, flags.Height, flags.Output)
}

func GenerateFrames(input string, frequency int, width int, height int) []string {
	var conf = configuration.Configuration{
		FfmpegPath:  "ffmpeg",
		FfprobePath: "ffprobe",
	}

	ffprobe := ffprobe.NewFfProbe(&conf, nil)
	probe, err := ffprobe.GetProbe(input)
	if err != nil {
		panic(fmt.Sprintf("could probe: %v", err))
	}

	duration, err := strconv.ParseFloat(probe.GetFirstVideoStream().Duration, 64)
	if err != nil {
		panic(fmt.Sprintf("could not get video duration: %v", err))
	}

	frames := make([]string, 0, 0)
	for i := 0; i < int(duration); i += frequency {
		frame := Extract(input, i, width, height)
		frames = append(frames, frame)
	}

	return frames
}

func Extract(input string, seek int, width int, height int) string {
	seekString := strconv.Itoa(seek)
	output := random_helper.Generate(10, random_helper.AZAndCaps) + ".png"
	cmd := exec.Command(
		"ffmpeg",
		"-ss",
		seekString,
		"-i",
		input,
		"-vf",
		fmt.Sprintf("scale=%dx%d", width, height),
		"-frames",
		"1",
		output,
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		panic(fmt.Sprintf("failed to extract frame: %v", err))
	}
	return output
}

func Montage(frames []string, output string) {
	cmd := exec.Command("montage", "-mode", "concatenate", "-tile", "10x10")

	for _, image := range frames {
		cmd.Args = append(cmd.Args, image)
	}

	cmd.Args = append(cmd.Args, output)
	_, err := cmd.CombinedOutput()
	Clean(frames)
	if err != nil {
		panic(fmt.Sprintf("failed to created montage: %v", err))
	}
}

func GenerateVtt(frames []string, frequency int, width int, height int, output *string) {

	file, err := os.Create(*output)
	if err != nil {
		panic("failed to create vtt file")
	}

	file.WriteString("WEBVTT\n")
	t1 := time.Second * 0

	for range frames {
		t2 := time.Duration(frequency) * time.Second
		start := time_helper.FormatHHMMSSmm(t1)
		end := time_helper.FormatHHMMSSmm(t1 + t2)

		line := fmt.Sprintf("%s --> %s %s#xywh=%d,%d,%d,%d\n", start, end, *output, x, y, width, height)
		file.WriteString(line)
		t1 = t2
	}

}

func Clean(frames []string) {
	for _, frame := range frames {
		os.Remove(frame)
	}
}
