package main

import (
	"flag"
	"github.com/giorgisio/goav/avdevice"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/pkg/cmd"
)

func main() {
	showFormats := flag.Bool("formats", false, "show available formats")
	showMuxers := flag.Bool("muxers", false, "show available muxers")
	showDemuxers := flag.Bool("demuxers", false, "show available demuxers")
	showDevices := flag.Bool("devices", false, "show available devices")
	showCodecs := flag.Bool("codecs", false, "show available codecs")
	showDecoders := flag.Bool("decoders", false, "show available decoders")
	showEncoders := flag.Bool("encoders", false, "show available encoders")
	showBSFS := flag.Bool("bsfs", false, "show available bit stream filters")
	showProtocols := flag.Bool("protocols", false, "show available protocols")
	showFilters := flag.Bool("filters", false, "show available filters")
	showPixFmts := flag.Bool("pix_fmts", false, "show available pixel formats")
	showLayouts := flag.Bool("layouts", false, "show standard channel layouts")
	showSampleFmts := flag.Bool("sample_fmts", false, "show available audio sample formats")
	showColors := flag.Bool("colors", false, "show available color names")
	flag.Parse()

	avformat.AvRegisterAll()
	avformat.AvformatNetworkInit()
	defer avformat.AvformatNetworkDeinit()
	avdevice.AvdeviceRegisterAll()

	cmd.ShowVersion()

	if *showFormats {
		cmd.ShowFormats()
		return
	} else if *showMuxers {
		cmd.ShowMuxers()
	} else if *showDemuxers {
		cmd.ShowDemuxers()
	} else if *showDevices {
		cmd.ShowDevices()
	} else if *showCodecs {
		cmd.ShowCodecs()
	} else if *showDecoders {
		cmd.ShowDecoders()
	} else if *showEncoders {
		cmd.ShowEncoders()
	} else if *showBSFS {
		cmd.ShowBitstreamFilters()
	} else if *showProtocols {
		cmd.ShowProtocols()
	} else if *showFilters {
		cmd.ShowFilters()
	} else if *showPixFmts {
		cmd.ShowPixelFormats()
	} else if *showLayouts {
		cmd.ShowChannelLayouts()
	} else if *showSampleFmts {
		cmd.ShowSampleFormats()
	} else if *showColors {
		cmd.ShowColors()
	}
}
