package main

import (
	"flag"
	"fmt"
	"github.com/Paxx-RnD/go-ffmpeg/configuration"
	"github.com/Paxx-RnD/go-ffmpeg/ffprobe"
	"sprite-preview/sprite"
	"sprite-preview/types"
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

	var conf = configuration.Configuration{
		FfmpegPath:  "ffmpeg",
		FfprobePath: "ffprobe",
	}

	ffprobe := ffprobe.NewFfProbe(&conf, nil)
	s := sprite.NewService(flags, ffprobe)
	fmt.Println("Generating frames")
	frames := s.GenerateFrames()
	fmt.Printf("Found %d frames\n", len(frames))
	fmt.Printf("Start Montage\n")
	s.Montage(frames)

	if flags.Vtt == "" && flags.Vtt == "" {
		return
	}

	fmt.Printf("Generate Vtt")
	s.GenerateVtt(frames)
}
