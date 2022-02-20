// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// Giorgis (habtom@giorgis.io)

package avfilter

/*
#cgo pkg-config: libavfilter
#include <libavfilter/avfilter.h>
static int has_process_command(const AVFilter *filter) {return filter->process_command != NULL;}
*/
import "C"
import "unsafe"

//Get a filter definition matching the given name.
func AvfilterGetByName(n string) *Filter {
	return (*Filter)(C.avfilter_get_by_name(C.CString(n)))
}

func AvfilterIterate(opaque *unsafe.Pointer) *Filter {
	return (*Filter)(C.av_filter_iterate(opaque))
}

func (f *Filter) Name() string {
	return C.GoString(f.name)
}

func (f *Filter) Description() string {
	return C.GoString(f.description)
}

func (f *Filter) Flags() int {
	return int(f.flags)
}

func (f *Filter) IsDynamicInputs() bool {
	return f.Flags() & C.AVFILTER_FLAG_DYNAMIC_INPUTS != 0
}

func (f *Filter) IsDynamicOutputs() bool {
	return f.Flags() & C.AVFILTER_FLAG_DYNAMIC_OUTPUTS != 0
}

func (f *Filter) IsTimeline() bool {
	return f.Flags() & C.AVFILTER_FLAG_SUPPORT_TIMELINE != 0
}

func (f *Filter) IsSliceThreads() bool {
	return f.Flags() & C.AVFILTER_FLAG_SLICE_THREADS != 0
}

func (f *Filter) Inputs() *Pad {
	return (*Pad)(f.inputs)
}

func (f *Filter) Outputs() *Pad {
	return (*Pad)(f.outputs)
}

func (f *Filter) HasProcessCommand() bool {
	return int(C.has_process_command((*C.struct_AVFilter)(f))) != 0
}

//Register a filter.
func (f *Filter) AvfilterRegister() int {
	return int(C.avfilter_register((*C.struct_AVFilter)(f)))
}

//Iterate over all registered filters.
func (f *Filter) AvfilterNext() *Filter {
	return (*Filter)(C.avfilter_next((*C.struct_AVFilter)(f)))
}
