package ffmpeg

import "github.com/leslie-wang/goav/pkg/cmd"

type Conf struct {
	*cmd.ShowCmdConf
	Format string
	CodecNames string
	Presets string
	OptMap string
	OptMapChannel string
	OptMapMetadata string
	Input string
	FileOverwrite bool
	NoFileOverwrite bool
	IgnoreUnknownStreams bool
	CopyUnknownStreams bool
}

