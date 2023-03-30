# Sprite Preview
sprite-preview is a tool that extracts frames from a video file and creates a sprite sheet of these frames. It also generates a VTT file that allows displaying the video using the sprite sheet.

# Installation
To use sprite-preview, you must have Go installed on your system. You can then install sprite-preview using the following command:

```bash
go get github.com/Paxx-RnD/sprite-preview
```

# Usage

```bash
sprite-preview -i input_file.mp4 -o output.png -f 1 -w 160 -h 90 -v output.vtt
```

# Flags
- i: Path to the input video file.
- o: Path to the output sprite sheet file. 
- col: Number of columns in the sprite sheet. 
- row: Number of rows in the sprite sheet.
- f: Frequency of frames extraction in seconds.
- w: Width of each frame in pixels.
- h: Height of each frame in pixels.
- vtt: Path to the output VTT file. If not specified, no VTT file will be generated.

# Preview
