package dotnet

/*
#cgo CXXFLAGS: -std=c++11 -Wall -pedantic
#cgo linux LDFLAGS: -ldl
#include <stdio.h>
#include <stdlib.h>
#include "binding.hpp"

*/
import "C"

import (
	"github.com/kardianos/osext"

	"errors"
	"strings"
	"unsafe"
	"os"

	"fmt"
)

type TheFunction C.TheFunction

type Runtime struct {
	Params RuntimeParams
}

// Hold CLR initialization parameters
type RuntimeParams struct {
	ExePath                     string
	AppDomainFriendlyName       string
	Properties                  map[string]string
	ManagedAssemblyAbsolutePath string

	CLRFilesAbsolutePath	string
}

type Callback struct {
	f *func()
}

var Callbacks map[int]Callback

const DefaultAppDomainFriendlyName string = "app"

// Creates a new runtime.
func NewRuntime(params RuntimeParams) (runtime Runtime, err error) {
	runtime = Runtime{Params: params}
	err = runtime.Init()

	return runtime, err
}

// Runtime initialization
// This function sets a few default values to make everything easier.
func (r *Runtime) Init() (err error) {
	if r.Params.ExePath == "" {
		r.Params.ExePath, err = osext.Executable()
	}

	if r.Params.AppDomainFriendlyName == "" {
		r.Params.AppDomainFriendlyName = DefaultAppDomainFriendlyName
	}

	// In case you don't set APP_PATHS/NATIVE_DLL_SEARCH_DIRECTORIES, the package assumes your assemblies are in the same directory.
	if r.Params.Properties["APP_PATHS"] == "" && r.Params.Properties["NATIVE_DLL_SEARCH_DIRECTORIES"] == "" {
		executableFolder, _ := osext.ExecutableFolder()
		r.Params.Properties["APP_PATHS"] = executableFolder
		r.Params.Properties["NATIVE_DLL_SEARCH_DIRECTORIES"] = executableFolder
	}

	propertyCount := len(r.Params.Properties)

	propertyKeys := make([]string, 0, len(r.Params.Properties))
	propertyValues := make([]string, 0, len(r.Params.Properties))

	for k, v := range r.Params.Properties {
		propertyKeys = append(propertyKeys, k)
		propertyValues = append(propertyValues, v)
	}

	ExePath := C.CString(r.Params.ExePath)
	AppDomainFriendlyName := C.CString(r.Params.AppDomainFriendlyName)
	PropertyCount := C.int(propertyCount)
	PropertyKeys := C.CString(strings.Join(propertyKeys, ";"))
	PropertyValues := C.CString(strings.Join(propertyValues, ";"))

	var CLRFilesAbsolutePath string

	// CLRCommonPaths holds possible SDK locations
	var CLRCommonPaths []string = []string{
		"/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
		"/usr/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
	}

	// Test for common SDK paths, return err if they don't exist
	if r.Params.CLRFilesAbsolutePath == "" {
		for _, p := range CLRCommonPaths {
			_, err := os.Stat(p)
			if err == nil {
				CLRFilesAbsolutePath = p
				break
			}
		}

		if CLRFilesAbsolutePath == "" {
			err = errors.New("No SDK found")
			return err
		}
	} else {
		CLRFilesAbsolutePath = r.Params.CLRFilesAbsolutePath
	}

	CLRFilesAbsolutePathC := C.CString(CLRFilesAbsolutePath)

	ManagedAssemblyAbsolutePath := C.CString(r.Params.ManagedAssemblyAbsolutePath)

	// Call the binding
	var result C.int
	result = C.initializeCoreCLR(ExePath, AppDomainFriendlyName, PropertyCount, PropertyKeys, PropertyValues, ManagedAssemblyAbsolutePath, CLRFilesAbsolutePathC)

	if result == -1 {
		err = errors.New("Runtime error")
	}

	C.free(unsafe.Pointer(ExePath))
	C.free(unsafe.Pointer(AppDomainFriendlyName))
	C.free(unsafe.Pointer(PropertyKeys))
	C.free(unsafe.Pointer(PropertyValues))
	C.free(unsafe.Pointer(ManagedAssemblyAbsolutePath))
	C.free(unsafe.Pointer(CLRFilesAbsolutePathC))

	return err
}

// Unloads the current app
// Reference:
//
//	https://github.com/dotnet/coreclr/blob/d81d773312dcae24d0b5d56cb972bf71e22f856c/src/dlls/mscoree/unixinterface.cpp#L281
//
func (r *Runtime) Shutdown() (err error) {
	var result C.int
	result = C.shutdownCoreCLR()

	if result == -1 {
		err = errors.New("Shutdown error.")
	}

	return err
}

// Loads an assembly file and call the default entrypoint.
func (r *Runtime) ExecuteManagedAssembly(assembly string) (err error) {
	var result C.int
	CAssembly := C.CString(assembly)
	result = C.executeManagedAssembly(CAssembly)
	C.free(unsafe.Pointer(CAssembly))

	if result == -1 {
		err = errors.New("Can't execute")
	}

	return err
}

// Makes it possible to call .NET stuff from Go.
func (r *Runtime) CreateDelegate(assemblyName string, typeName string, methodName string) func() {

	// var err error

	CAssemblyName := C.CString(assemblyName)
	CTypeName := C.CString(typeName)
	CMethodName := C.CString(methodName)

	var result C.int

	if result != 0 {
		// err = errors.New("Can't create delegate")
	}

	return func() {
		result = C.createDelegate(CAssemblyName, CTypeName, CMethodName, 1)
	}
}

func RegisterCallback(f func()) int {
	fmt.Println("Registering callback!", len(Callbacks))
	var n = len(Callbacks)
	// Callbacks[n] = Callback{&f}
	callback := Callback{&f}

	if Callbacks == nil {
		Callbacks = make(map[int]Callback)
	}

	Callbacks[n] = callback
	return n
}
