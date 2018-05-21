package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

// Generator generates code for interoperability between Go and .NET.
type Generator struct {
	PkgName string
	Input   []*Input
}

// Input is a data structure used for every parsed file, contains the AST, path information and annotations.
type Input struct {
	fileset *token.FileSet
	astFile *ast.File
	path    string

	annotations []Annotation

	cgoCode        *bytes.Buffer
	goCode         *bytes.Buffer
	bindingHeaders *bytes.Buffer
	bindingSource  *bytes.Buffer

	generator *Generator
}

// New initializes a generator.
func New(input []string) *Generator {
	g := Generator{
		Input: make([]*Input, 0),
	}
	for _, path := range input {
		i := &Input{
			path:           path,
			generator:      &g,
			cgoCode:        &bytes.Buffer{},
			goCode:         &bytes.Buffer{},
			bindingHeaders: &bytes.Buffer{},
			bindingSource:  &bytes.Buffer{},
		}
		g.Input = append(g.Input, i)
	}
	return &g
}

// Parse builds the AST and extracts information useful for the generation step.
func (i *Input) Parse() (err error) {
	log.WithFields(logrus.Fields{
		"file": i.path,
	}).Info("Parsing file")

	i.fileset = token.NewFileSet()
	i.astFile, err = parser.ParseFile(i.fileset, i.path, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// Extract the package name:
	i.generator.PkgName = i.astFile.Name.Name
	if i.generator.PkgName == "main" {
		err := errors.New("Use a non-main package name")
		log.Fatal(err)
	}
	log.Debugf("Found package \"%s\"", i.generator.PkgName)

	i.annotations = make([]Annotation, 0)
	// Walk through the declarations and function comments:
	for _, decl := range i.astFile.Decls {
		function, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if function.Doc == nil {
			log.Warnf("Skipping function '%s'", function.Name.Name)
			continue
		}
		comments := function.Doc.List
		if len(comments) == 0 {
			continue
		}
		for _, v := range comments {
			if strings.HasPrefix(v.Text, "// ") {
				annotation := i.parseFuncAnnotation(v.Text, function)
				log.WithFields(logrus.Fields{
					"annotation": annotation,
				}).Debug("Found annotation")
				i.annotations = append(i.annotations, annotation)
			}
		}
	}
	log.Debugf("Finished checking %s, found %d annotations", i.path, len(i.annotations))
	return nil
}

// Verbose modifies the verbosity level.
func (g *Generator) Verbose(mode bool) *Generator {
	if mode {
		log.Level = logrus.DebugLevel
	}
	return g
}

// Parse iterates through every input file, calling the parse method.
func (g *Generator) Parse() (err error) {
	for _, input := range g.Input {
		err = input.Parse()
		if err != nil {
			return err
		}
	}
	return nil
}

// Generate performs the code generation step.
func (g *Generator) Generate() (err error) {
	// TODO: find out when/how to handle multiple files
	mainFile := g.Input[0]
	if mainFile == nil {
	}

	mainFile.cgoCode.WriteString("/*\n#include <binding.hpp>")

	for _, a := range mainFile.annotations {
		a.Render()
	}

	mainFile.cgoCode.WriteString("\n*/\n")
	mainFile.cgoCode.WriteString("import \"C\"")

	goCode := &bytes.Buffer{}
	goCode.WriteString("package " + g.PkgName + "\n")
	mainFile.cgoCode.WriteTo(goCode)
	mainFile.goCode.WriteTo(goCode)

	renderedHeaders := bytes.Buffer{}
	bindingsData := map[string]interface{}{
		"HeaderDefinitions": mainFile.bindingHeaders.String(),
	}
	bindingsTemplate := template.Must(template.New("binding.hpp").ParseFiles("/Users/matias/.gvm/pkgsets/go1.7/global/src/github.com/matiasinsaurralde/go-dotnet/generator/templates/binding.hpp"))
	err = bindingsTemplate.Execute(&renderedHeaders, &bindingsData)

	renderedImpl := bytes.Buffer{}
	implData := map[string]interface{}{
		"Impls": mainFile.bindingSource.String(),
	}
	implTemplate := template.Must(template.New("binding.cpp").ParseFiles("/Users/matias/.gvm/pkgsets/go1.7/global/src/github.com/matiasinsaurralde/go-dotnet/generator/templates/binding.cpp"))
	err = implTemplate.Execute(&renderedImpl, &implData)

	originalPath := "/Users/matias/.gvm/pkgsets/go1.7/global/src/github.com/matiasinsaurralde/go-dotnet/dotnet"

	basePath := "/Users/matias/.gvm/pkgsets/go1.7/global/src/github.com/matiasinsaurralde/go-dotnet/pkg"
	os.Mkdir(basePath, 0755)
	copyFiles, _ := ioutil.ReadDir(originalPath)
	for _, v := range copyFiles {
		path := filepath.Join(originalPath, v.Name())
		destPath := filepath.Join(basePath, v.Name())
		fmt.Printf("Copy %s to %s\n", v.Name(), destPath)
		data, _ := ioutil.ReadFile(path)
		ioutil.WriteFile(destPath, data, 0755)
	}

	bindingsHeadersPath := filepath.Join(basePath, "binding.hpp")
	ioutil.WriteFile(bindingsHeadersPath, renderedHeaders.Bytes(), 0755)

	bindingsImplPath := filepath.Join(basePath, "binding.cpp")
	ioutil.WriteFile(bindingsImplPath, renderedImpl.Bytes(), 0755)

	pkgPath := filepath.Join(basePath, "pkg.go")
	ioutil.WriteFile(pkgPath, goCode.Bytes(), 0755)

	return nil
}
