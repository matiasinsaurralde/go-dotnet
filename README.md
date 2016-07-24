# go-dotnet

This is a PoC Go wrapper for the .NET Core Runtime, this project uses ```cgo``` and has been tested under OSX. It covers two basic use cases that are provided by the [CLR Hosting API](https://blogs.msdn.microsoft.com/msdnforum/2010/07/09/use-clr4-hosting-api-to-invoke-net-assembly-from-native-c/):

* Load and run a .exe, using its default entrypoint, just like [corerun](https://github.com/dotnet/coreclr/blob/master/src/coreclr/hosts/unixcorerun/corerun.cpp) and [coreconsole](https://github.com/dotnet/coreclr/blob/master/src/coreclr/hosts/unixcoreconsole/coreconsole.cpp) do, check ```ExecuteManagedAssembly```.

* Load a .dll, setup [delegates](http://www.fancy-development.net/hosting-net-core-clr-in-your-own-process) and call them from your Go functions.

## An example

```
package main

import (
	"github.com/matiasinsaurralde/go-dotnet"

	"fmt"
	"os"
)

func main() {
	fmt.Println("This is main, I'll initialize the .NET runtime.")

	properties := map[string]string{
		"TRUSTED_PLATFORM_ASSEMBLIES":   "",
        "APP_PATHS":                     "/Users/matias/dev/dotnet/cdotnet/lib/HelloWorld",
        "NATIVE_DLL_SEARCH_DIRECTORIES": "/Users/matias/dev/dotnet/cdotnet/lib/HelloWorld:/usr/local/share/dotnet/shared/Microsoft.NETCore.App/1.0.0",
	}

	err, runtime := dotnet.NewRuntime(dotnet.RuntimeParams{
		Properties:                  properties,
	})

	if err != nil {
		fmt.Println("Something bad happened! :(")
		os.Exit(1)
	}

	fmt.Println("Runtime loaded.")

	SayHello := runtime.CreateDelegate("HelloWorld", "HelloWorld.HelloWorld", "Hello")

    // this will call HelloWorld.HelloWorld.Hello :)
	SayHello()

	runtime.Shutdown()
}
```

## Ideas

* Run some benchmarks.
* Add/enhance ```net/http``` samples, like [this](https://github.com/matiasinsaurralde/go-dotnet/blob/master/examples/http.go).

## License

[MIT](https://github.com/matiasinsaurralde/go-dotnet/blob/master/LICENSE)
