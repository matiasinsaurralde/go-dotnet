package dotnet

/*
#include <stdlib.h>
#include "binding.hpp"
*/
import "C"

import (
	// "fmt"
	"unsafe"
)

type DelegateArguments struct {
	p unsafe.Pointer
}

func (d *DelegateArguments) Size() int {
	var sz C.int
	sz = C.getSize(d.p)
	return int(sz)
}

func (d *DelegateArguments) Append(i interface{}) {
	var argType int
	var argValue unsafe.Pointer

	switch i.(type) {
	case int:
		argType = 0
		value := i.(int)
		var CValue C.int
		CValue = C.int(value)
		argValue = unsafe.Pointer(&CValue)
	}

	if argValue != nil {
		C.pushDelegateArgument(d.p, argValue, C.int(argType))
	}
}

func NewDelegateArguments(args []interface{}) DelegateArguments {
	d := DelegateArguments{}
	d.p = C.initDelegateArguments()

	for _, v := range args {
		d.Append(v)
	}

	return d
}
