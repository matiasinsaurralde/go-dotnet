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
  "strings"
  "fmt"
)

type Runtime struct {
  Params RuntimeParams
}

type RuntimeParams struct {
  ExePath string
  AppDomainFriendlyName string
  Properties map[string]string
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

  fmt.Println("ExePath = ", r.Params.ExePath)
  fmt.Println("AppDomainFriendlyName = ", r.Params.AppDomainFriendlyName)
  fmt.Println("PropertyCount = ", propertyCount)

  propertyKeys := make([]string, 0, len(r.Params.Properties))
  propertyValues := make([]string, 0, len(r.Params.Properties))

  for k, v := range r.Params.Properties {
    propertyKeys = append(propertyKeys, k)
    propertyValues = append(propertyValues, v)
  }

  ExePath := C.CString(r.Params.ExePath)
  AppDomainFriendlyName := C.CString(r.Params.AppDomainFriendlyName)
  PropertyCount := C.int(propertyCount)
  PropertyKeys := C.CString(strings.Join(propertyKeys, ":"))
  PropertyValues := C.CString(strings.Join(propertyValues, ":"))

  C.initializeCoreCLR(ExePath, AppDomainFriendlyName, PropertyCount, PropertyKeys, PropertyValues)

  C.free(unsafe.Pointer(ExePath))
  C.free(unsafe.Pointer(AppDomainFriendlyName))
  C.free(unsafe.Pointer(PropertyKeys))
  C.free(unsafe.Pointer(PropertyValues))

  return err
}
