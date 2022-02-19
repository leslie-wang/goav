package main

import (
	"github.com/giorgisio/goav/avdevice"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/pkg/cmd"
)

func main() {
	avformat.AvRegisterAll()
	avdevice.AvdeviceRegisterAll()
	ctx := avformat.AvformatAllocContext()
	cmd.ShowFormats(ctx)
}
