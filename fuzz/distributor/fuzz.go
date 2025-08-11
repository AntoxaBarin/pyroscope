//go:build libfuzzer
// +build libfuzzer

package main

import (
	"fmt"
	profilev1 "github.com/grafana/pyroscope/api/gen/proto/go/google/v1"
	"unsafe"

	mutator "github.com/yandex-cloud/go-protobuf-mutator"
)

// #include <stdint.h>
import "C"

//export LLVMFuzzerInitialize
func LLVMFuzzerInitialize(argc *C.int, argv ***C.char) C.int {

	//todo fuzz limits as well
	return 0
}

//export LLVMFuzzerCustomMutator
func LLVMFuzzerCustomMutator(data *C.char, size C.size_t, maxSize C.size_t, seed C.uint) C.size_t {
	gdata := unsafe.Slice((*byte)(unsafe.Pointer(data)), size)
	message := new(profilev1.Profile)
	if err := message.UnmarshalVT(gdata); err != nil {
		return 0
	}
	mutator := mutator.New(int64(seed), int(maxSize-size))
	if err := mutator.MutateProto(message); err != nil {
		fmt.Printf("Failed to mutate message: %+v", err)
		return 0
	}
	gdata = unsafe.Slice((*byte)(unsafe.Pointer(data)), maxSize)
	if message.SizeVT() > len(gdata) {
		return 0
	}

	if sz, err := message.MarshalToVT(gdata); err != nil {
		panic(err)
	} else {
		return C.size_t(sz)
	}
}

//export LLVMFuzzerTestOneInput
func LLVMFuzzerTestOneInput(data *C.char, size C.size_t) C.int {
	gdata := unsafe.Slice((*byte)(unsafe.Pointer(data)), size)
	message := new(profilev1.Profile)
	if err := message.UnmarshalVT(gdata); err != nil {
		//fmt.Printf("Failed to parse proto message: %+v", err)
		return 0
	}
	panic("unmarshalled")
	//todo multiple distributor configurations

	return 0
}

func main() {

}
