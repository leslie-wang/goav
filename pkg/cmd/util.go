package cmd

/*
#cgo pkg-config: libavutil

#include <libavformat/avformat.h>
#include <libavutil/avutil.h>

static int is_device(const AVClass *avclass)
{
    if (!avclass)
        return 0;
		return AV_IS_INPUT_DEVICE(avclass->category) || AV_IS_OUTPUT_DEVICE(avclass->category);
}
*/
import "C"

import (
	"fmt"

	"github.com/giorgisio/goav/avformat"
)

type ShowMuxDemux int

const (
	ShowDefault ShowMuxDemux = iota
	ShowDemuxers
	ShowMuxers
)

func ShowFormats(ctx *avformat.Context) {
	showFormatsDevice(ctx, false, ShowDefault)
}

func showFormatsDevice(ctx *avformat.Context, deviceOnly bool, showMuxDemux ShowMuxDemux) {
	prefix := "Devices:"
	if !deviceOnly {
		prefix = "File formats:"
	}
	fmt.Printf("%s\n D. = Demuxing supported\n .E = Muxing supported\n --\n", prefix)

	lastName := ""
	ofmt := ctx.Oformat()
	ifmt := ctx.Iformat()
	for {
		decode := " "
		encode := " "
		name := ""
		longName := ""
		if showMuxDemux != ShowDemuxers {
			for {
				ofmt = ofmt.AvOformatNext()
				if ofmt == nil {
					break
				}
				if !ofmt.IsDevice() && deviceOnly {
					fmt.Printf("--- mux device: %s\n", ofmt.Name())
					continue
				}
				if (name == "" || ofmt.Name() < name) && ofmt.Name() > lastName {
					name = ofmt.Name()
					longName = ofmt.LongName()
					encode = "E"
				}
			}
		}
		if showMuxDemux != ShowMuxers {
			for {
				ifmt = ifmt.AvIformatNext()
				if ifmt == nil {
					break
				}
				if !ifmt.IsDevice() && deviceOnly {
					fmt.Printf("--- demux device: %s\n", ifmt.Name())
					continue
				}
				if (name == "" || ifmt.Name() < name) && ifmt.Name() > lastName {
					name = ifmt.Name()
					longName = ifmt.LongName()
					encode = " "
				}
				if name == ifmt.Name() {
					decode = "D"
				}
			}
		}
		if name == "" {
			break
		}
		lastName = name
		fmt.Printf(" %s%s %-15s %s\n", decode, encode, name, longName)
	}
}
