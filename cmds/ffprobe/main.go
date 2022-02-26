package main

import (
	"flag"
	"github.com/leslie-wang/goav/avdevice"
	"github.com/leslie-wang/goav/avformat"
	"github.com/leslie-wang/goav/pkg/cmd"
	"log"
	"os"
)
type conf struct {
	*cmd.ShowCmdConf
	forceFmt string
}

func main() {
	conf := &conf{ShowCmdConf: &cmd.ShowCmdConf{}}
	f := flag.NewFlagSet("ffmpeg", flag.ExitOnError)
	conf.SetFlags(f)

	avformat.AvRegisterAll()
	avformat.AvformatNetworkInit()
	defer avformat.AvformatNetworkDeinit()
	avdevice.AvdeviceRegisterAll()

	cmd.ShowVersion()

	if len(os.Args) == 1 {
		f.Usage()
		return
	}
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// return if it is only show command
	if conf.Show() {
		return
	}
}
