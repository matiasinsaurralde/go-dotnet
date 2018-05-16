package mybinding

/*
#cgo CXXFLAGS: -std=c++11 -Wall
#cgo linux LDFLAGS: -ldl
#include <stdlib.h>
#include "runtime.hpp"
*/
import "C"

import (
	"github.com/kardianos/osext"

	"errors"
	"os"
	"strings"
	"unsafe"
)

var runtimeInstance = &Runtime{}

// Runtime is the runtime data structure.
type Runtime struct {
	Params        RuntimeParams
	delegateSetup func()
}

// RuntimeParams holds the CLR initialization parameters
type RuntimeParams struct {
	ExePath                     string
	AppDomainFriendlyName       string
	Properties                  map[string]string
	ManagedAssemblyAbsolutePath string

	CLRFilesAbsolutePath string
}

const DefaultAppDomainFriendlyName string = "app"

// Init performs the runtime initialization
// This function sets a few default values to make everything easier.
func Init() (err error) {
	if runtimeInstance.Params.ExePath == "" {
		runtimeInstance.Params.ExePath, err = osext.Executable()
	}

	if runtimeInstance.Params.AppDomainFriendlyName == "" {
		runtimeInstance.Params.AppDomainFriendlyName = DefaultAppDomainFriendlyName
	}

	if runtimeInstance.Params.Properties == nil {
		runtimeInstance.Params.Properties = make(map[string]string)
	}

	// In case you don't set APP_PATHS/NATIVE_DLL_SEARCH_DIRECTORIES, the package assumes your assemblies are in the same directory.
	if runtimeInstance.Params.Properties["APP_PATHS"] == "" && runtimeInstance.Params.Properties["NATIVE_DLL_SEARCH_DIRECTORIES"] == "" {
		executableFolder, _ := osext.ExecutableFolder()
		runtimeInstance.Params.Properties["APP_PATHS"] = executableFolder
		runtimeInstance.Params.Properties["NATIVE_DLL_SEARCH_DIRECTORIES"] = executableFolder
	}

	propertyCount := len(runtimeInstance.Params.Properties)

	propertyKeys := make([]string, 0, len(runtimeInstance.Params.Properties))
	propertyValues := make([]string, 0, len(runtimeInstance.Params.Properties))

	for k, v := range runtimeInstance.Params.Properties {
		propertyKeys = append(propertyKeys, k)
		propertyValues = append(propertyValues, v)
	}

	ExePath := C.CString(runtimeInstance.Params.ExePath)
	AppDomainFriendlyName := C.CString(runtimeInstance.Params.AppDomainFriendlyName)
	PropertyCount := C.int(propertyCount)
	PropertyKeys := C.CString(strings.Join(propertyKeys, ";"))
	PropertyValues := C.CString(strings.Join(propertyValues, ";"))

	var CLRFilesAbsolutePath string

	// CLRCommonPaths holds possible SDK locations
	var CLRCommonPaths = []string{
		"/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
		"/usr/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
	}

	// Test for common SDK paths, return err if they don't exist
	if runtimeInstance.Params.CLRFilesAbsolutePath == "" {
		for _, p := range CLRCommonPaths {
			_, err = os.Stat(p)
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
		CLRFilesAbsolutePath = runtimeInstance.Params.CLRFilesAbsolutePath
	}

	CLRFilesAbsolutePathC := C.CString(CLRFilesAbsolutePath)

	ManagedAssemblyAbsolutePath := C.CString(runtimeInstance.Params.ManagedAssemblyAbsolutePath)

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

	runtimeInstance.delegateSetup()

	return err
}

// Shutdown unloads the current app
//
//	https://github.com/dotnet/coreclr/blob/d81d773312dcae24d0b5d56cb972bf71e22f856c/src/dlls/mscoree/unixinterface.cpp#L281
//
func (r *Runtime) Shutdown() (err error) {
	var result C.int
	result = C.shutdownCoreCLR()

	if result == -1 {
		err = errors.New("Shutdown error")
	}

	return err
}

// CreateDelegate wraps a cgo call to coreclr_create_delegate, receives a function pointer.
func CreateDelegate(assemblyName string, typeName string, methodName string, delegateID int, f *unsafe.Pointer) int {
	CassemblyName := C.CString(assemblyName)
	CtypeName := C.CString(typeName)
	CmethodName := C.CString(methodName)
	CdelegateID := C.int(0)
	result := C.createDelegate(CassemblyName, CtypeName, CmethodName, CdelegateID, f)
	return int(result)
}

// SetupDelegates sets all create_delegate calls to be executed after the runtime initialization.
func SetupDelegates(f func()) {
	runtimeInstance.delegateSetup = f
}
