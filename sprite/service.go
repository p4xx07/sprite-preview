package sprite

import (
	"fmt"
	"github.com/Paxx-RnD/go-ffmpeg/ffprobe"
	"github.com/Paxx-RnD/go-helper/helpers/math_helper/mathint"
	"github.com/Paxx-RnD/go-helper/helpers/path_helper"
	"github.com/Paxx-RnD/go-helper/helpers/random_helper"
	"github.com/Paxx-RnD/go-helper/helpers/time_helper"
	"math"
	"os"
	"os/exec"
	"path"
	"sprite-preview/types"
	"strconv"
	"time"
)

type IService interface {
	GenerateFrames() []string
	Montage(frames []string)
	GenerateVtt(frames []string)
}

type service struct {
	flags   types.Flags
	ffprobe ffprobe.IFfprobe
}

func NewService(flags types.Flags, ffprobe ffprobe.IFfprobe) IService {
	return &service{flags: flags, ffprobe: ffprobe}
}

func (s *service) GenerateFrames() []string {
	probe, err := s.ffprobe.GetProbe(*s.flags.Input)
	if err != nil {
		panic(fmt.Sprintf("could probe: %v", err))
	}

	duration, err := strconv.ParseFloat(probe.GetFirstVideoStream().Duration, 64)
	if err != nil {
		panic(fmt.Sprintf("could not get video duration: %v", err))
	}

	frames := make([]string, 0, 0)
	for i := 0; i < int(duration); i += *s.flags.Frequency {
		frame := s.extract(i)
		frames = append(frames, frame)
	}

	return frames
}

func (s *service) extract(seek int) string {
	seekString := strconv.Itoa(seek)
	output := random_helper.Generate(10, random_helper.AZAndCaps) + ".png"
	cmd := exec.Command(
		"ffmpeg",
		"-ss",
		seekString,
		"-i",
		*s.flags.Input,
		"-vf",
		fmt.Sprintf("scale=%dx%d", *s.flags.Width, *s.flags.Height),
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

func (s *service) Montage(frames []string) {
	cmd := exec.Command("montage", "-mode", "concatenate", "-tile", "10x10")

	for _, image := range frames {
		cmd.Args = append(cmd.Args, image)
	}

	cmd.Args = append(cmd.Args, *s.flags.Output)
	_, err := cmd.CombinedOutput()
	s.clean(frames)
	if err != nil {
		panic(fmt.Sprintf("failed to created montage: %v", err))
	}
}

func (s *service) GenerateVtt(frames []string) {
	file, err := os.Create(*s.flags.Vtt)
	if err != nil {
		panic("failed to create vtt file")
	}

	file.WriteString("WEBVTT\n\n")
	t1 := time.Second * 0

	grid := *s.flags.Columns * *s.flags.Rows
	max := float64(len(frames)) / float64(grid)
	nSprites := mathint.Max(1, int(math.Ceil(max)))
	for n := 0; n < nSprites; n++ {
		output := fmt.Sprintf("%s-%d%s", path_helper.BasePathWithoutExt(*s.flags.Output), n, path.Ext(*s.flags.Output))
		for y := 0; y < *s.flags.Columns; y++ {
			for x := 0; x < *s.flags.Rows; x++ {
				t2 := time.Duration(*s.flags.Frequency) * time.Second
				start := time_helper.FormatHHMMSSmm(t1)
				end := time_helper.FormatHHMMSSmm(t1 + t2)

				line := fmt.Sprintf("%s --> %s %s#xywh=%d,%d,%d,%d\n\n", start, end, output, x**s.flags.Width, y**s.flags.Height, *s.flags.Width, *s.flags.Height)
				file.WriteString(line)
				t1 += t2
			}
		}
	}

}

func (s *service) clean(frames []string) {
	for _, frame := range frames {
		os.Remove(frame)
	}
}
