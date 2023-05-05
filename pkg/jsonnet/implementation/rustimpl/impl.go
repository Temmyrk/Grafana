package rustimpl

// #cgo CFLAGS: -I /Users/julienduchesne/Repos/jrsonnet
// #cgo LDFLAGS: -L /Users/julienduchesne/Repos/jrsonnet/target/x86_64-unknown-linux-gnu/release/ -ljsonnet
// #include "bindings/c/libjsonnet.h"
import "C"

import (
	"errors"

	"github.com/grafana/tanka/pkg/jsonnet/implementation/types"
)

type JsonnetRustVM struct {
	vm *C.struct_JsonnetVm
}

func (vm *JsonnetRustVM) EvaluateAnonymousSnippet(filename, snippet string) (string, error) {
	errorInt := C.int(0)
	result := C.jsonnet_evaluate_snippet(vm.vm, C.CString(filename), C.CString(snippet), &errorInt)
	if errorInt != 0 {
		return "", errors.New(C.GoString(result))
	}
	return C.GoString(result), nil
}

func (vm *JsonnetRustVM) EvaluateFile(filename string) (string, error) {
	errorInt := C.int(0)
	result := C.jsonnet_evaluate_file(vm.vm, C.CString(filename), &errorInt)
	if errorInt != 0 {
		return "", errors.New(C.GoString(result))
	}
	return C.GoString(result), nil
}

type JsonnetRustImplementation struct{}

func (i *JsonnetRustImplementation) MakeVM(importPaths []string, extCode map[string]string, tlaCode map[string]string, maxStack int) types.JsonnetVM {
	return &JsonnetRustVM{
		vm: makeVM(importPaths, extCode, tlaCode, maxStack),
	}
}

func makeVM(importPaths []string, extCode map[string]string, tlaCode map[string]string, maxStack int) *C.struct_JsonnetVm {
	vm := C.jsonnet_make()
	for _, path := range importPaths {
		C.jsonnet_jpath_add(vm, C.CString(path))
	}
	for key, value := range extCode {
		C.jsonnet_ext_code(vm, C.CString(key), C.CString(value))
	}
	for key, value := range tlaCode {
		C.jsonnet_tla_code(vm, C.CString(key), C.CString(value))
	}
	if maxStack > 0 {
		C.jsonnet_max_stack(vm, C.uint(maxStack))
	}
	return vm
}
