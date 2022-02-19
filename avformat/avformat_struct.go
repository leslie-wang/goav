package avformat

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

func (of OutputFormat) Name() string {
	return C.GoString(of.name)
}

func (of OutputFormat) LongName() string {
	return C.GoString(of.long_name)
}

func (of OutputFormat) PrivClass() *C.struct_AVClass {
	return of.priv_class
}

func (of OutputFormat) IsDevice() bool {
	return C.is_device(of.priv_class) != 0
}

func (inf InputFormat) Name() string {
	return C.GoString(inf.name)
}

func (inf InputFormat) LongName() string {
	return C.GoString(inf.long_name)
}

func (inf InputFormat) PrivClass() *C.struct_AVClass {
	return inf.priv_class
}

func (inf InputFormat) IsDevice() bool {
	return C.is_device(inf.priv_class) != 0
}
