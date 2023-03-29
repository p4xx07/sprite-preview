package types

import (
	"flag"
	"fmt"
)

type Flags struct {
	Input       *string
	Output      *string
	GenerateVtt bool
	Frequency   int
	Width       int
	Height      int
	Help        bool
}

func (f *Flags) Set() error {
	f.Input = flag.String("i", "", "path of the input video")
	f.Output = flag.String("o", "", "output sprites")
	f.Frequency = *flag.Int("f", 3, "extract frames every n seconds")
	f.Width = *flag.Int("w", 160, "frame width")
	f.Height = *flag.Int("h", 90, "frame height")
	f.GenerateVtt = *flag.Bool("vtt", false, "specify the output for the vtt file")
	f.Help = *flag.Bool("help", false, "shows help")

	flag.Parse()

	if f.Input == nil {
		return fmt.Errorf("need to specify input")
	}
	if f.Output == nil {
		return fmt.Errorf("need to specify output")
	}
	if f.Frequency <= 0 {
		return fmt.Errorf("need to specify valid target")
	}

	return nil
}
