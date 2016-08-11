package main

import(
  "go/parser"
  "go/token"
  "go/ast"

  "fmt"
)

func main() {
  fset := token.NewFileSet()
  astf, _ := parser.ParseFile(fset, "playground/sample.go", nil, parser.ParseComments)

  for _, decl := range astf.Decls {
    var functionDeclaration *ast.FuncDecl
    functionDeclaration = decl.(*ast.FuncDecl)

    var returnFields, recvFields []*ast.Field
    returnFields = functionDeclaration.Type.Results.List
    recvFields = functionDeclaration.Type.Params.List

    fmt.Println("Function Declaration:")
    fmt.Println(" Name:", functionDeclaration.Name)

    fmt.Println(" Receive fields:" )
    for _, f := range recvFields {
      fmt.Println( "    ", f )
    }

    fmt.Println(" Return Fields:" )
    for _, f := range returnFields {
      fmt.Println(  "    ", f )
    }
  }
}
