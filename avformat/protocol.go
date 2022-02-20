package avformat

//#cgo pkg-config: libavformat
//#include <libavformat/avformat.h>
import "C"
import "unsafe"

func AvioEnumProtocols (opaque *unsafe.Pointer, output bool) string {
	op := 0
	if output {
		op = 1
	}
	return C.GoString(C.avio_enum_protocols(opaque, C.int(op)))
}
