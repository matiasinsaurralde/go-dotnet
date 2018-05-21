package dotnet

/*
#cgo CXXFLAGS: -std=c++11 -Wall
#cgo linux LDFLAGS: -ldl
#include <stdlib.h>
#include "runtime.hpp"
*/
import "C"

import (
	"errors"
	"os"
	"strings"
	"unsafe"

	"github.com/kardianos/osext"
)

const (
	assemblyNotFound       = 0x80070002
	typeLoadException      = 0x80131522
	missingMethodException = 0x80131513
	nullReferenceException = 0x80004003

	defaultAppDomainFriendlyName = "app"
)

var (
	runtimeInstance = &Runtime{}

	errAssemblyNotFound       = errors.New("Assembly not found")
	errTypeLoadException      = errors.New("Missing type")
	errMissingMethodException = errors.New("Missing method")
	errNullReferenceException = errors.New("Invalid delegate function pointer")
)

// Runtime is the runtime data structure.
type Runtime struct {
	Params        RuntimeParams
	delegateSetup func() error
}

// RuntimeParams holds the CLR initialization parameters
type RuntimeParams struct {
	ExePath                     string
	AppDomainFriendlyName       string
	Properties                  map[string]string
	ManagedAssemblyAbsolutePath string

	CLRFilesAbsolutePath string
}

// SetParams sets initial runtime parameters.
func SetParams(params RuntimeParams) {
	runtimeInstance.Params = params
}

// Init performs the runtime initialization
// This function sets a few default values to make everything easier.
func Init() (err error) {
	if runtimeInstance.Params.ExePath == "" {
		runtimeInstance.Params.ExePath, err = osext.Executable()
	}

	if runtimeInstance.Params.AppDomainFriendlyName == "" {
		runtimeInstance.Params.AppDomainFriendlyName = defaultAppDomainFriendlyName
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

	count := len(runtimeInstance.Params.Properties)

	keys := make([]string, 0, len(runtimeInstance.Params.Properties))
	vals := make([]string, 0, len(runtimeInstance.Params.Properties))

	for k, v := range runtimeInstance.Params.Properties {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	exePath := C.CString(runtimeInstance.Params.ExePath)
	appDomainFriendlyName := C.CString(runtimeInstance.Params.AppDomainFriendlyName)
	propertyCount := C.int(count)
	propertyKeys := C.CString(strings.Join(keys, ";"))
	propertyValues := C.CString(strings.Join(vals, ";"))

	var clrFilesAbsolutePath string

	// clrCommonPaths holds possible SDK locations
	var clrCommonPaths = []string{
		"/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
		"/usr/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
	}

	// Test for common SDK paths, return err if they don't exist
	if runtimeInstance.Params.CLRFilesAbsolutePath == "" {
		for _, p := range clrCommonPaths {
			_, err = os.Stat(p)
			if err == nil {
				clrFilesAbsolutePath = p
				break
			}
		}

		if clrFilesAbsolutePath == "" {
			err = errors.New("No SDK found")
			return err
		}
	} else {
		clrFilesAbsolutePath = runtimeInstance.Params.CLRFilesAbsolutePath
	}

	clrFilesAbsolutePathC := C.CString(clrFilesAbsolutePath)

	managedAssemblyAbsolutePath := C.CString(runtimeInstance.Params.ManagedAssemblyAbsolutePath)

	// Call the binding
	var result C.int
	result = C.initializeCoreCLR(exePath, appDomainFriendlyName, propertyCount, propertyKeys, propertyValues, managedAssemblyAbsolutePath, clrFilesAbsolutePathC)

	if result == -1 {
		err = errors.New("Runtime error")
	}

	C.free(unsafe.Pointer(exePath))
	C.free(unsafe.Pointer(appDomainFriendlyName))
	C.free(unsafe.Pointer(propertyKeys))
	C.free(unsafe.Pointer(propertyValues))
	C.free(unsafe.Pointer(managedAssemblyAbsolutePath))
	C.free(unsafe.Pointer(clrFilesAbsolutePathC))

	// No delegates set?
	if runtimeInstance.delegateSetup == nil {
		return nil
	}
	return runtimeInstance.delegateSetup()
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
func CreateDelegate(assembly string, typ string, method string, delegate int, f *unsafe.Pointer) error {
	assemblyName := C.CString(assembly)
	typeName := C.CString(typ)
	methodName := C.CString(method)
	delegateID := C.int(delegate)
	result := C.createDelegate(assemblyName, typeName, methodName, delegateID, f)
	code := uint32(result)
	switch code {
	case assemblyNotFound:
		return errAssemblyNotFound
	case typeLoadException:
		return errTypeLoadException
	case missingMethodException:
		return errMissingMethodException
	case nullReferenceException:
		return errNullReferenceException
	}
	return nil
}

// SetupDelegates sets all create_delegate calls to be executed after the runtime initialization.
func SetupDelegates(f func() error) {
	runtimeInstance.delegateSetup = f
}
