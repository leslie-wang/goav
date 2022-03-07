package ffmpeg

import (
	"github.com/pkg/errors"
	"github.com/leslie-wang/goav/avcodec"
	"github.com/leslie-wang/goav/avformat"
	"github.com/leslie-wang/goav/avutil"
)

func findCodec(name string, typ avutil.MediaType, isEncoder bool) (*avcodec.Codec, error) {
	var codec *avcodec.Codec
	codecStr := ""
	if isEncoder {
		codecStr = "encoder"
		codec = avcodec.AvcodecFindEncoderByName(name)
	} else {
		codecStr = "decoder"
		codec = avcodec.AvcodecFindDecoderByName(name)
	}
	if codec == nil {
		desc := avcodec.AvcodecDescriptorGetByName(name)
		if desc != nil {
			if isEncoder {
				codec = avcodec.AvcodecFindEncoder(desc.ID())
			} else {
				codec = avcodec.AvcodecFindDecoder(desc.ID())
			}
		}
	}
	if codec == nil {
		return nil, errors.Errorf("Unknown %s '%s'", codecStr, name)
	}
	if codec.Type() != typ {
		return nil, errors.Errorf("Invalid %s type '%s'", codecStr, name)
	}
	return codec, nil
}

func chooseDecoder(name string, stream *avformat.Stream) (*avcodec.Codec, error) {
	if name != "" {
		codec, err := findCodec(name, stream.CodecParameters().MediaType(), false)
		if err != nil {
			return nil, err
		}
		stream.CodecParameters().SetID(codec.ID())
		return codec, nil
	}
	codec := avcodec.AvcodecFindDecoder(stream.CodecParameters().ID())
	if codec == nil {
		return nil, errors.Errorf("Unknown decoder '%s'", name)
	}
	return codec, nil
}

func chooseEncoder(name string, conf *OutputConfig, ctx *avformat.Context, stream *avformat.Stream) error {
	typ := stream.CodecParameters().MediaType()
	switch typ {
	case avutil.AVMEDIA_TYPE_VIDEO:
	case avutil.AVMEDIA_TYPE_AUDIO:
	case avutil.AVMEDIA_TYPE_SUBTITLE:
		if name == "" {
			oname := ""
			id := avformat.AvGuessCodec(ctx.Oformat(), "", ctx.URL(), "", typ)
			stream.CodecParameters().SetID(id)

			conf.Enc = avcodec.AvcodecFindEncoder(id)
			if conf.Enc == nil {
				return errors.Wrapf(avutil.ErrorFromCode(avutil.AvErrorEncoderNotFound),
					"Automatic encoder selection failed for output stream #0:0. " +
					"Default encoder for format %s (codec %s) is " +
				"probably disabled. Please choose an encoder manually",
					oname, avcodec.AvcodecGetName(id))
			}
		} else if name == CodecCopy {
			conf.StreamCopy = true
		} else {
			var err error
			conf.Enc, err = findCodec(name, typ, true)
			if err != nil {
				return err
			}
			stream.CodecParameters().SetID(conf.Enc.ID())
		}
	default:
		conf.StreamCopy = true
	}
	return nil
}
