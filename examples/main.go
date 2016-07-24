package main

import (
	"github.com/matiasinsaurralde/go-dotnet"

	"fmt"
	"os"
)

func main() {
	fmt.Println("This is main, I'll initialize the .NET runtime.")

	err, runtime := dotnet.NewRuntime()

	if err != nil {
		fmt.Println("Something bad happened! :(")
		os.Exit(1)
	}

	fmt.Println("Runtime loaded:", runtime)
}
