// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

// Package avutil is a utility library to aid portable multimedia programming.
// It contains safe portable string functions, random number generators, data structures,
// additional mathematics functions, cryptography and multimedia related functionality.
// Some generic features and utilities provided by the libavutil library
package avutil

//#cgo pkg-config: libavutil
//#include <libavutil/pixdesc.h>
//#include <libavutil/channel_layout.h>
//#include <libavutil/samplefmt.h>
//#include <libavutil/parseutils.h>
//#include <libavutil/avutil.h>
//#include <stdlib.h>
import "C"
import (
	"unsafe"
)

type (
	Options       C.struct_AVOptions
	AvTree        C.struct_AVTree
	Rational      C.struct_AVRational
	MediaType     C.enum_AVMediaType
	AvPictureType C.enum_AVPictureType
	PixelFormat   C.enum_AVPixelFormat
	PixFmtDescr   C.struct_AVPixFmtDescriptor
	SampleFormat  C.enum_AVSampleFormat
	File          C.FILE
)

func AVVersionMajor(v uint) uint {
	return v >> 16
}

func AVVersionMinor(v uint) uint {
	return (v & 0x00FF00) >> 8
}

func AVVersionMicro(v uint) uint {
	return v & 0xFF
}

//Return the LIBAvUTIL_VERSION_INT constant.
func AvutilVersion() (uint, uint, uint) {
	v := uint(C.avutil_version())
	return AVVersionMajor(v), AVVersionMinor(v), AVVersionMicro(v)
}

//Return the libavutil build-time configuration.
func AvutilConfiguration() string {
	return C.GoString(C.avutil_configuration())
}

//Return the libavutil license.
func AvutilLicense() string {
	return C.GoString(C.avutil_license())
}

//Return a string describing the media_type enum, NULL if media_type is unknown.
func AvGetMediaTypeString(mt MediaType) string {
	return C.GoString(C.av_get_media_type_string((C.enum_AVMediaType)(mt)))
}

//Return a single letter to describe the given picture type pict_type.
func AvGetPictureTypeChar(pt AvPictureType) string {
	return string(C.av_get_picture_type_char((C.enum_AVPictureType)(pt)))
}

//Return x default pointer in case p is NULL.
func AvXIfNull(p, x int) {
	C.av_x_if_null(unsafe.Pointer(&p), unsafe.Pointer(&x))
}

//Compute the length of an integer list.
func AvIntListLengthForSize(e uint, l int, t uint64) uint {
	return uint(C.av_int_list_length_for_size(C.uint(e), unsafe.Pointer(&l), (C.uint64_t)(t)))
}

//Open a file using a UTF-8 filename.
func AvFopenUtf8(p, m string) *File {
	f := C.av_fopen_utf8(C.CString(p), C.CString(m))
	return (*File)(f)
}

//Return the fractional representation of the internal time base.
func AvGetTimeBaseQ() Rational {
	return (Rational)(C.av_get_time_base_q())
}

func (m MediaType) ToChar() byte {
	switch m {
	case C.AVMEDIA_TYPE_VIDEO:    return 'V'
	case C.AVMEDIA_TYPE_AUDIO:    return 'A'
	case C.AVMEDIA_TYPE_DATA:     return 'D'
	case C.AVMEDIA_TYPE_SUBTITLE: return 'S'
	case C.AVMEDIA_TYPE_ATTACHMENT:return 'T'
	default:                    return '?'
	}
}

func (pfd *PixFmtDescr) AvPixFmtDescNext() *PixFmtDescr{
	return (*PixFmtDescr)(C.av_pix_fmt_desc_next((*C.struct_AVPixFmtDescriptor)(pfd)))
}

func (pfd *PixFmtDescr) Name() string {
	return C.GoString(pfd.name)
}

func (pfd *PixFmtDescr) NbComponents() byte {
	return byte(pfd.nb_components)
}

func (pfd *PixFmtDescr) Flags() int {
	return int(pfd.flags)
}

func (pfd *PixFmtDescr) IsHWAccel() bool {
	return pfd.Flags() & C.AV_PIX_FMT_FLAG_HWACCEL != 0
}

func (pfd *PixFmtDescr) IsPal() bool {
	return pfd.Flags() & C.AV_PIX_FMT_FLAG_PAL != 0
}

func (pfd *PixFmtDescr) IsBitstream() bool {
	return pfd.Flags() & C.AV_PIX_FMT_FLAG_BITSTREAM != 0
}

func (pfd *PixFmtDescr) PixelFormat() PixelFormat {
	return (PixelFormat)(C.av_pix_fmt_desc_get_id((*C.struct_AVPixFmtDescriptor)(pfd)))
}

func (pfd *PixFmtDescr) BitsPerPixel() int {
	return int(C.av_get_bits_per_pixel((*C.struct_AVPixFmtDescriptor)(pfd)))
}

func GetChannelName(index uint64) string {
	return C.GoString(C.av_get_channel_name((C.uint64_t)(index)))
}

func GetChannelDescription(index uint64) string {
	return C.GoString(C.av_get_channel_description((C.uint64_t)(index)))
}

func GetStandardChannelLayout(index uint) (int, uint64, string) {
	var layout uint64
	var name *C.char
	defer C.free(unsafe.Pointer(name))

	ret := C.av_get_standard_channel_layout((C.uint)(index), (*C.uint64_t)(&layout), &name)
	return int(ret), layout, C.GoString(name)
}

func GetSampleFmtString(sampleFmt SampleFormat) string {
	size := 128
	descr := (*C.char)(C.malloc(C.ulong(size)))
	defer C.free(unsafe.Pointer(descr))

	C.av_get_sample_fmt_string(descr, C.int(size), C.enum_AVSampleFormat(sampleFmt))
	return C.GoString(descr)
}

func GetKnownColorName(colorIdx int) (string, []byte) {
	var rgb *C.uint8_t
	name := C.GoString(C.av_get_known_color_name(C.int(colorIdx), (**C.uint8_t)(&rgb)))
	if name == "" {
		return "", nil
	}
	return name, C.GoBytes(unsafe.Pointer(rgb), 3)
}
