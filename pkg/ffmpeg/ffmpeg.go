package ffmpeg

import (
	"github.com/leslie-wang/goav/avcodec"
	"github.com/leslie-wang/goav/avformat"
	"github.com/leslie-wang/goav/avutil"
	"github.com/leslie-wang/goav/pkg/cmd"
	"github.com/pkg/errors"
	"syscall"
)

type InputConfig struct {
	StartTime int64
	StopTime int64
	RecordingTime int64
	Format string
	FrameRate float32
	FormatOpts *avutil.Dictionary
	start int64
}

type Input struct {
	ctx *avformat.Context
	conf *InputConfig
	filename string
}

func NewInput(conf *InputConfig, filename string) (*Input, error) {
	var iformat *avformat.InputFormat
	if conf.Format != "" {
		iformat = avformat.AvFindInputFormat(conf.Format)
		if iformat == nil {
			return nil, errors.Errorf("Unknown input format: '%s'", conf.Format)
		}
	}

	ctx := avformat.AvformatAllocContext()
	if ctx == nil {
		return nil, errors.Errorf("%s: %s", filename, syscall.ENOMEM.Error())
	}
	if conf.FormatOpts == nil {
		conf.FormatOpts = avutil.NewDict()
	}

	scanAllPmtsSet := false
	entry := conf.FormatOpts.AvDictGet(OptScannAllFmts, nil, avutil.AV_DICT_MATCH_CASE)
	if entry == nil {
		val := "1"
		conf.FormatOpts.AvDictSet(OptScannAllFmts, &val, avutil.AV_DICT_DONT_OVERWRITE)
		scanAllPmtsSet = true
	}
	ret := avformat.AvformatOpenInput(&ctx, filename, iformat, &conf.FormatOpts)
	if ret < 0 {
		return nil, errors.Wrapf(avutil.ErrorFromCode(ret), filename)
	}

	if scanAllPmtsSet {
		conf.FormatOpts.AvDictSet(OptScannAllFmts, nil, avutil.AV_DICT_MATCH_CASE)
	}

	for _, stream := range ctx.Streams() {
		chooseDecoder("", stream)
	}
	input := &Input{ctx: ctx, conf: conf, filename: filename}
	return input, nil
}

func addInputStreams(ctx *avformat.Context) {

}

type OutputConfig struct {
	FrameRate float32
	FrameNumber int
	StreamCopy bool
	Enc *avcodec.Codec
	EncCtx *avcodec.Context
}

type Output struct {
	ctx *avformat.Context
	conf *OutputConfig
	filename string
}

func newOutputVideoStream() {

}

func newOutputAudioStream() {

}

func newOutputStream(ctx *avformat.Context, typ avutil.MediaType) (*avformat.Stream, error){
	stream := ctx.AvformatNewStream(nil)
	if stream == nil {
		return nil, errors.New("Could not alloc stream.")
	}
	stream.CodecParameters().SetMediaType(typ)

	return stream, nil
}

func NewOutput(conf *OutputConfig, filename string) (*Output, error) {
	ctx, err := avformat.AvformatAllocOutputContext2(nil, nil, filename)
	if err != nil {
		return nil, errors.Wrapf(err, filename)
	}
	output := &Output{ctx: ctx, conf: conf, filename: filename}
	return output, nil
}

type Conf struct {
	*cmd.ShowCmdConf
	Format string
	CodecNames string
	Presets string
	OptMap string
	OptMapChannel string
	OptMapMetadata string
	FileOverwrite bool
	NoFileOverwrite bool
	IgnoreUnknownStreams bool
	CopyUnknownStreams bool
	Input []*Input
	Output []*Output
}

func initInputStream(conf *Conf) error {
	return nil
}

func transcodeInit(conf *Conf) error {
	return nil
}

func Transcode(conf *Conf) error {
	err := transcodeInit(conf)
	if err != nil {
		return err
	}
	return nil
}