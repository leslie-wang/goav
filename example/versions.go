package main

import (
	"log"

	"github.com/leslie-wang/goav/avcodec"
	"github.com/leslie-wang/goav/avdevice"
	"github.com/leslie-wang/goav/avfilter"
	"github.com/leslie-wang/goav/avformat"
	"github.com/leslie-wang/goav/avutil"
	"github.com/leslie-wang/goav/swresample"
	"github.com/leslie-wang/goav/swscale"
)

func main() {

	// Register all formats and codecs
	avformat.AvRegisterAll()
	avcodec.AvcodecRegisterAll()

	log.Printf("AvFilter Version:\t%v", avfilter.AvfilterVersion())
	log.Printf("AvDevice Version:\t%v", avdevice.AvdeviceVersion())
	log.Printf("SWScale Version:\t%v", swscale.SwscaleVersion())
	log.Printf("AvUtil Version:\t%v", avutil.AvutilVersion())
	log.Printf("AvCodec Version:\t%v", avcodec.AvcodecVersion())
	log.Printf("Resample Version:\t%v", swresample.SwresampleLicense())

}
