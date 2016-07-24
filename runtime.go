package dotnet

/*
#include <stdio.h>
#include <stdlib.h>
#include "binding.hpp"

static void abc() {
};

*/
import "C"

import(
  "github.com/kardianos/osext"

  "unsafe"
  "errors"
  "strings"
)

type Runtime struct {
  Params RuntimeParams
}

type RuntimeParams struct {
  ExePath string
  AppDomainFriendlyName string
  Properties map[string]string
  ManagedAssemblyAbsolutePath string
}

const DefaultAppDomainFriendlyName string = "app"

func NewRuntime(params RuntimeParams) (err error, runtime Runtime) {
  runtime = Runtime{Params: params}
  // C.initializeCoreCLR()
  // C.executeAssembly()

  err = runtime.Init()

  C.abc()

  return err, runtime
}

func( r *Runtime ) Init() (err error) {
  if r.Params.ExePath == "" {
    r.Params.ExePath, err = osext.Executable()
  }
  if r.Params.AppDomainFriendlyName == "" {
    r.Params.AppDomainFriendlyName = DefaultAppDomainFriendlyName
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

  CLRFilesAbsolutePath := C.CString("/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0")

  ManagedAssemblyAbsolutePath := C.CString(r.Params.ManagedAssemblyAbsolutePath)

  var result C.int
  result = C.initializeCoreCLR(ExePath, AppDomainFriendlyName, PropertyCount, PropertyKeys, PropertyValues, ManagedAssemblyAbsolutePath, CLRFilesAbsolutePath)

  if result == -1 {
    err = errors.New("Runtime error")
  }

  C.free(unsafe.Pointer(ExePath))
  C.free(unsafe.Pointer(AppDomainFriendlyName))
  C.free(unsafe.Pointer(PropertyKeys))
  C.free(unsafe.Pointer(PropertyValues))
  C.free(unsafe.Pointer(ManagedAssemblyAbsolutePath))
  C.free(unsafe.Pointer(CLRFilesAbsolutePath))

  return err
}

func( r *Runtime ) Shutdown() (err error) {
  var result C.int
  result = C.shutdownCoreCLR()

  if result == -1 {
    err = errors.New("Shutdown error.")
  }

  return err
}

func( r *Runtime) ExecuteManagedAssembly(assembly string) (err error) {
  var result C.int
  CAssembly := C.CString(assembly)
  result = C.executeManagedAssembly(CAssembly)
  C.free(unsafe.Pointer(CAssembly))

  if result == -1 {
    err = errors.New("Can't execute")
  }

  return err
}

func( r *Runtime) CreateDelegate() {
  C.createDelegate()
}
