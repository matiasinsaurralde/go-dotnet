## Code generation

At the moment, the interaction with CoreCLR delegates is pretty simple and makes it possible to call .NET methods from Go.

However I'm missing **two important things**:

* It's not possible to pass arguments to a delegate.
* It's not possible to access returned values.

I've been trying to modify the package to do both things and keep the original proposed syntax.

My proposed syntax looks like this:
```go
SayHello := CreateDelegateCallback("HelloWorld.Hello")
SayHello("Rob Pike", func(newValue string) {
  fmt.Println("Delegate returns:", newValue)
})
```

On the .NET side, the `Hello` function will accept and return a string. Calling `Hello("Rob Pike")` will return `"Hello Rob Pike"`.

I think that it should be possible to take return values directly too:

```go
SayHelloAgain := CreateDelegateReturn("HelloWorld.Hello")
SayHelloAgainReturns := SayHelloAgain("Bill Gates")

fmt.Println("Delegate returns:", SayHelloAgainReturns)
```

This is very cool stuff and I think that it will be safer and easier to have a code generation tool (like Protocol Buffers do with `go generate`), that builds the Go bindings for your .NET methods (and .NET data structures, why not?).

The Go code to match these prospective syntax looks like this:
```go
func CreateDelegateCallback( delegateName string ) func( string, func(string) ) {
  f := func( inputValue string, callback func(string) ) {
    newValue := fmt.Sprintf( "Hello %s", inputValue )
    callback(newValue)
  }
  return f
}

func CreateDelegateReturn( delegateName string ) func(args... interface{}) interface{} {
  f := func(args... interface{}) interface{} {
    inputValue := args[0]
    newValue := fmt.Sprintf( "Hello %s", inputValue)
    return newValue
  }
  return f
}
```

After delegate setup, we'll need to make a heavy use of `interface{}` and type checks, when passing arguments or retrieving return values from .NET. So I've been thinking that a code generator could make this better, so you write some functions and prefix them with line comments containing the required delegate information (namespace/name, arguments, return data type).

```go
package mybinding

//create_delegate: HelloWorld.Hello(string) string
func SayHello(string) string {
}
```

Then you run `go generate` and all the Hosted API magic stuff is ready to import and use.
