package main

import (
	"github.com/matiasinsaurralde/go-dotnet"

	"fmt"
	"os"
)

func main() {
	fmt.Println("This is main, I'll initialize the .NET runtime.")

  properties := map[string]string{
    "TRUSTED_PLATFORM_ASSEMBLIES": "trusted",
    "APP_PATHS": "apppaths",
    "APP_NI_PATHS": "appnipaths",
    "NATIVE_DLL_SEARCH_DIRECTORIES": "d",
    "AppDomainCompatSwitch": "e",
  }

	err, runtime := dotnet.NewRuntime(dotnet.RuntimeParams{
    Properties: properties,
  })

  fmt.Println( "RuntimeParams = ", runtime.Params )

	if err != nil {
		fmt.Println("Something bad happened! :(")
		os.Exit(1)
	}

	fmt.Println("Runtime loaded:", runtime)
}
