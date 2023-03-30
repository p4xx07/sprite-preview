package main

import (
	"flag"
	"github.com/Paxx-RnD/go-ffmpeg/configuration"
	"github.com/Paxx-RnD/go-ffmpeg/ffprobe"
	"peak-tracer/sprite"
	"peak-tracer/types"
)

func main() {
	flags := types.Flags{}
	err := flags.Set()
	if err != nil {
		panic(err)
	}

	if *flags.Help {
		flag.Usage()
		return
	}

	var conf = configuration.Configuration{
		FfmpegPath:  "ffmpeg",
		FfprobePath: "ffprobe",
	}

	ffprobe := ffprobe.NewFfProbe(&conf, nil)
	s := sprite.NewService(flags, ffprobe)
	frames := s.GenerateFrames()
	s.Montage(frames)

	if flags.Vtt == nil && *flags.Vtt == "" {
		return
	}

	s.GenerateVtt(frames)
}
