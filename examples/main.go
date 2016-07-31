package main

import (
	"github.com/matiasinsaurralde/go-dotnet"

	"fmt"
	"os"
)

func main() {
	fmt.Println("This is main, I'll initialize the .NET runtime.")

	/*
		If you don't set the TRUSTED_PLATFORM_ASSEMBLIES, it will use the default tpaList value.
	*/

	properties := map[string]string{
		// "TRUSTED_PLATFORM_ASSEMBLIES": "/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0/mscorlib.ni.dll:/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0/System.Private.CoreLib.ni.dll",
		"APP_PATHS":                     "/Users/matias/dev/dotnet/cdotnet/lib/HelloWorld",
		"NATIVE_DLL_SEARCH_DIRECTORIES": "/Users/matias/dev/dotnet/cdotnet/lib/HelloWorld:/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
	}

	runtime, err := dotnet.NewRuntime(dotnet.RuntimeParams{
		Properties:                  properties,
	})
	defer runtime.Shutdown()

	if err != nil {
		fmt.Println("Something bad happened! :(")
		os.Exit(1)
	}

	fmt.Println("Runtime loaded.")

	// runtime.ExecuteManagedAssembly("HelloWorldMain.exe")

	SayHello := runtime.CreateDelegate("HelloWorld", "HelloWorld.HelloWorld", "Hello")
	SayHello()
}
