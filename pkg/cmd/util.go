package cmd

/*
#cgo pkg-config: libavutil

#include <libavutil/avutil.h>
#include <libavformat/avformat.h>
#include <libavfilter/avfilter.h>

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
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avdevice"
	"github.com/giorgisio/goav/avfilter"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swresample"
	"github.com/giorgisio/goav/swscale"
	"sort"
	"strings"
	"unsafe"

	"github.com/giorgisio/goav/avformat"
)

func ShowVersion() {
	major, minor, micro := avutil.AvutilVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "avutil", major, minor, micro)
	major, minor, micro = avcodec.AvcodecVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "avcodec", major, minor, micro)
	major, minor, micro = avformat.AvformatVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "avformat", major, minor, micro)
	major, minor, micro = avdevice.AvdeviceVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "avdevice", major, minor, micro)
	major, minor, micro = avfilter.AvfilterVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "avfilter", major, minor, micro)
	major, minor, micro = swresample.SwresampleVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "swresample", major, minor, micro)
	major, minor, micro = swscale.SwscaleVersion()
	fmt.Printf("  lib%-11s %2d.%3d.%3d\n", "swscale", major, minor, micro)
}

type ShowMuxDemux int

const (
	showDefault ShowMuxDemux = iota
	showDemuxers
	showMuxers
)

func ShowFormats() {
	showFormatsDevice(false, true, true)
}

func ShowMuxers() {
	showFormatsDevice(false, true, false)
}

func ShowDemuxers() {
	showFormatsDevice(false, false, true)
}

func ShowDevices() {
	showFormatsDevice(true, true, true)
}

func showFormatsDevice(deviceOnly, showMux, showDemux bool) {
	prefix := "Devices:"
	if !deviceOnly {
		prefix = "File formats:"
	}
	fmt.Printf("%s\n D. = Demuxing supported\n .E = Muxing supported\n --\n", prefix)

	ctx := avformat.AvformatAllocContext()
	defer ctx.AvformatFreeContext()
	ofmt := ctx.Oformat()
	ifmt := ctx.Iformat()
	lastName := ""
	for {
		decode := " "
		encode := " "
		name := ""
		longName := ""
		if showMux {
			for {
				ofmt = ofmt.AvOformatNext()
				if ofmt == nil {
					break
				}
				if !ofmt.IsDevice() && deviceOnly {
					continue
				}
				if (name == "" || ofmt.Name() < name) && ofmt.Name() > lastName {
					name = ofmt.Name()
					longName = ofmt.LongName()
					encode = "E"
				}
			}
		}
		if showDemux {
			for {
				ifmt = ifmt.AvIformatNext()
				if ifmt == nil {
					break
				}
				if !ifmt.IsDevice() && deviceOnly {
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

func findCodecs(decoder bool, id avcodec.CodecId) []*avcodec.Codec{
	var codecs []*avcodec.Codec
	p := &avcodec.Codec{}
	p = p.AvCodecNext() // now it should start from nil
	for {
		p = p.AvCodecNext()
		if p == nil {
			break
		}
		if p.ID() != id {
			continue
		}
		ok := false
		if decoder {
			ok = p.AvCodecIsDecoder()
		} else {
			ok = p.AvCodecIsEncoder()
		}
		if ok {
			codecs = append(codecs, p)
		}
	}
	return codecs
}


func listCodecs() []*avcodec.Descriptor {
	desc := &avcodec.Descriptor{}
	desc = desc.AvcodecDescriptorNext()
	codecs := []*avcodec.Descriptor{}
	for {
		desc = desc.AvcodecDescriptorNext()
		if desc == nil {
			break
		}
		if strings.Contains(desc.Name(), "_deprecated") {
			continue
		}
		codecs = append(codecs, desc)
	}
	sort.Slice(codecs, func(i, j int) bool {
		x := codecs[i].Type()
		y := codecs[j].Type()
		if x != y {
			return x < y
		}
		return codecs[i].Name() < codecs[j].Name()
	})
	return codecs
}

func ShowCodecs() {
	fmt.Printf("Codecs:\n" +
	" D..... = Decoding supported\n" +
	" .E.... = Encoding supported\n" +
	" ..V... = Video codec\n" +
	" ..A... = Audio codec\n" +
	" ..S... = Subtitle codec\n" +
	" ...I.. = Intra frame-only codec\n" +
	" ....L. = Lossy compression\n" +
	" .....S = Lossless compression\n" +
	" -------\n")

	codecs := listCodecs()
	for _, c := range codecs {
		fmt.Printf(" ")
		if avcodec.AvcodecFindDecoder(c.ID()) != nil {
			fmt.Printf("D")
		} else {
			fmt.Printf(".")
		}
		if avcodec.AvcodecFindEncoder(c.ID()) != nil {
			fmt.Printf("E")
		} else {
			fmt.Printf(".")
		}
		fmt.Printf("%c", c.Type().ToChar())
		if c.IsIntraOnly() {
			fmt.Printf("I")
		} else {
			fmt.Printf(".")
		}
		if c.IsLossy() {
			fmt.Printf("L")
		} else {
			fmt.Printf(".")
		}
		if c.IsLossless() {
			fmt.Printf("S")
		} else {
			fmt.Printf(".")
		}
		fmt.Printf(" %-20s %s", c.Name(), c.LongName())

		/* print decoders/encoders when there's more than one or their
		 * names are different from codec name */
		filterCodecs := findCodecs(true, c.ID())
		if len(filterCodecs) > 0 && (len(filterCodecs) != 1 || filterCodecs[0].Name() != c.Name()) {
			names := []string{}
			for _, fc := range filterCodecs {
				names = append(names, fc.Name())
			}
			fmt.Printf(" (decoders: %s )", strings.Join(names, " "))
		}
		filterCodecs = findCodecs(false, c.ID())
		if len(filterCodecs) > 0 && (len(filterCodecs) != 1 || filterCodecs[0].Name() != c.Name()) {
			names := []string{}
			for _, fc := range filterCodecs {
				names = append(names, fc.Name())
			}
			fmt.Printf(" (encoders: %s )", strings.Join(names, " "))
		}

		fmt.Printf("\n")
	}
}

func ShowEncoders() {
	showDeEncoders(false)
}

func ShowDecoders() {
	showDeEncoders(true)
}

func showDeEncoders(decoder bool) {
	prefix := "Encoders"
	if decoder {
		prefix = "Decoders"
	}
	fmt.Printf("%s:\n" +
	" V..... = Video\n" +
	" A..... = Audio\n" +
	" S..... = Subtitle\n" +
	" .F.... = Frame-level multithreading\n" +
	" ..S... = Slice-level multithreading\n" +
	" ...X.. = Codec is experimental\n" +
	" ....B. = Supports draw_horiz_band\n" +
	" .....D = Supports direct rendering method 1\n" +
	" ------\n", prefix)

	codecs := listCodecs()
	for _, c := range codecs {
		filterCodecs := findCodecs(decoder, c.ID())
		for _, fc := range filterCodecs {
			fmt.Printf(" %c", fc.Type().ToChar())
			if fc.IsCapFrameThreads() {
				fmt.Printf("F")
			} else {
				fmt.Printf(".")
			}
			if fc.IsCapSliceThreads() {
				fmt.Printf("S")
			} else {
				fmt.Printf(".")
			}
			if fc.IsCapExperimental() {
				fmt.Printf("X")
			} else {
				fmt.Printf(".")
			}
			if fc.IsCapDrawHorizBand() {
				fmt.Printf("B")
			} else {
				fmt.Printf(".")
			}
			if fc.IsCapDR1() {
				fmt.Printf("D")
			} else {
				fmt.Printf(".")
			}

			fmt.Printf(" %-20s %s", fc.Name(), fc.LongName())
			if fc.Name() != c.Name() {
				fmt.Printf(" (codec %s)", c.Name())
			}
			fmt.Printf("\n")
		}
	}
}

func ShowBitstreamFilters() {
	var (
		bsf *avcodec.BitStreamFilter
		opaque unsafe.Pointer
	)
	fmt.Println("Bitstream filters")
	for {
		bsf = avcodec.AvBitstreamFilterIterate(&opaque)
		if bsf == nil {
			break
		}
		fmt.Println(bsf.Name())
	}
	fmt.Println()
}

func ShowProtocols() {
	fmt.Println("Supported file protocols:")
	fmt.Println("Input:")

	var opaque unsafe.Pointer
	for {
		name := avformat.AvioEnumProtocols(&opaque, false)
		if name == "" {
			break
		}
		fmt.Printf("  %s\n", name)
	}

	fmt.Println("Output:")
	for {
		name := avformat.AvioEnumProtocols(&opaque, true)
		if name == "" {
			break
		}
		fmt.Printf("  %s\n", name)
	}
}

func ShowFilters() {
	fmt.Printf("Filters:\n" +
	"  T.. = Timeline support\n" +
	"  .S. = Slice threading\n" +
	"  ..C = Command support\n" +
	"  A = Audio input/output\n" +
	"  V = Video input/output\n" +
	"  N = Dynamic number and/or type of input/output\n" +
	"  | = Source or sink filter\n")

	var opaque unsafe.Pointer
	for {
		filter := avfilter.AvfilterIterate(&opaque)
		if filter == nil {
			break
		}
		descr := ""
		for i := 0; i < 2; i++ {
			var (
				j int
				pad *avfilter.Pad
			)
			if i == 0 {
				pad = filter.Inputs()
			} else {
				descr += "->"
				pad = filter.Outputs()
			}
			for {
				if pad == nil {
					break
				}
				name := avfilter.AvfilterPadGetName(pad, j)
				if name == "" {
					break
				}
				// hardcoded length from ffprobe
				if len(descr) >= 60 {
					break
				}
				descr += string(avfilter.AvfilterPadGetType(pad, j).ToChar())
				j++
			}
			if j == 0 {
				dy := "|"
				if (i == 0 && filter.IsDynamicInputs()) || (i != 0 && filter.IsDynamicOutputs()) {
					dy = "N"
				}
				descr += dy
			}
		}
		timeline := '.'
		if filter.IsTimeline() {
			timeline = 'T'
		}
		sliceThreads := '.'
		if filter.IsSliceThreads() {
			sliceThreads = 'S'
		}
		hasProcessCommand := '.'
		if filter.HasProcessCommand() {
			hasProcessCommand = 'C'
		}
		fmt.Printf(" %c%c%c %-17s %-10s %s\n", timeline, sliceThreads, hasProcessCommand,
			filter.Name(), descr, filter.Description())
	}
}

func ShowPixelFormats() {
	fmt.Printf("Pixel formats:\n" +
	"I.... = Supported Input  format for conversion\n" +
	".O... = Supported Output format for conversion\n" +
	"..H.. = Hardware accelerated format\n" +
	"...P. = Paletted format\n" +
	"....B = Bitstream format\n" +
	"FLAGS NAME            NB_COMPONENTS BITS_PER_PIXEL\n" +
	"-----\n")

	descr := &avutil.PixFmtDescr{}
	descr = descr.AvPixFmtDescNext()
	for {
		descr = descr.AvPixFmtDescNext()
		if descr == nil {
			break
		}
		pf := descr.PixelFormat()
		swsInput := '.'
		if swscale.SwsIssupportedinput(pf) {
			swsInput = 'I'
		}
		swsOutput := '.'
		if swscale.SwsIssupportedoutput(pf) {
			swsOutput = 'O'
		}
		hwAccel := '.'
		if descr.IsHWAccel() {
			hwAccel = 'H'
		}
		pal := '.'
		if descr.IsPal() {
			pal = 'P'
		}
		bs := '.'
		if descr.IsBitstream() {
			bs = 'B'
		}
		fmt.Printf("%c%c%c%c%c %-16s       %d            %2d\n", swsInput, swsOutput, hwAccel, pal, bs,
			descr.Name(), descr.NbComponents(), descr.BitsPerPixel())
	}
}

func ShowChannelLayouts() {
	fmt.Println("Individual channels:")
	fmt.Println("NAME           DESCRIPTION")
	for i := 0; i < 63; i++ {
		name := avutil.GetChannelName(1 << i)
		if name == "" {
			continue
		}
		fmt.Printf("%-14s %s\n", name, avutil.GetChannelDescription(1 << i))
	}
	fmt.Println("\nStandard channel layouts:")
	fmt.Println("NAME           DECOMPOSITION")
	var i uint
	for {
		ret, layout, name := avutil.GetStandardChannelLayout(i)
		if ret != 0 {
			break
		}
		if name != "" {
			fmt.Printf("%-14s ", name)
			var j uint64
			for j = 1; j != 0; j <<= 1 {
				if layout & j != 0 {
					if layout & (j - 1) != 0 {
						fmt.Printf("%s%s", "+", avutil.GetChannelName(j))
					} else {
						fmt.Printf("%s%s", "", avutil.GetChannelName(j))
					}
				}
			}
			fmt.Println()
		}
		i++
	}
}

func ShowSampleFormats() {
	for i := C.AV_SAMPLE_FMT_NONE; i < C.AV_SAMPLE_FMT_NB; i++ {
		fmt.Println(avutil.GetSampleFmtString(avutil.SampleFormat(i)))
	}
}

func ShowColors() {
	fmt.Printf("%-32s #RRGGBB\n", "name")

	i := 0
	for {
		name, rgb := avutil.GetKnownColorName(i)
		if name == "" {
			break
		}
		fmt.Printf("%-32s #%02x%02x%02x\n", name, rgb[0], rgb[1], rgb[2])
		i++
	}
}