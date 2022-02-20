// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

// Package swresample provides a high-level interface to the libswresample library audio resampling utilities
// The process of changing the sampling rate of a discrete signal to obtain a new discrete representation of the underlying continuous signal.
package swresample

/*
	#cgo pkg-config: libswresample
	#include <libswresample/swresample.h>
*/
import "C"
import "github.com/giorgisio/goav/avutil"

type (
	Context        C.struct_SwrContext
	Frame          C.struct_AVFrame
	Class          C.struct_AVClass
	AvSampleFormat C.enum_AVSampleFormat
)

//Get the Class for Context.
func SwrGetClass() *Class {
	return (*Class)(C.swr_get_class())
}

//Context constructor functions.Allocate Context.
func SwrAlloc() *Context {
	return (*Context)(C.swr_alloc())
}

//Configuration accessors
func SwresampleVersion() (uint, uint, uint) {
	v := uint(C.swresample_version())
	return avutil.AVVersionMajor(v), avutil.AVVersionMinor(v), avutil.AVVersionMicro(v)
}

func SwresampleConfiguration() string {
	return C.GoString(C.swresample_configuration())
}

func SwresampleLicense() string {
	return C.GoString(C.swresample_license())
}
