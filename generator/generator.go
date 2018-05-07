package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

// Generator generates code for interoperability between Go and .NET.
type Generator struct {
	Input []*Input
}

// Input is a data structure used for every parsed file, contains the AST, path information and annotations.
type Input struct {
	fileset *token.FileSet
	astFile *ast.File
	path    string

	annotations []*DelegateAnnotation
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
	i.annotations = make([]*DelegateAnnotation, 0)

	for _, decl := range i.astFile.Decls {
		function, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		comments := function.Doc.List
		if len(comments) == 0 {
			continue
		}
		for _, v := range comments {
			if strings.HasPrefix(v.Text, delegatePrefix) {
				annotation := parseAnnotation(v.Text)
				log.WithFields(logrus.Fields{
					"annotation": annotation,
				}).Debug("Found annotation")
				i.annotations = append(i.annotations, &annotation)
			}
		}
	}
	log.Debugf("Finished checking %s, found %d annotations", i.path, len(i.annotations))
	return nil
}

// New initializes a generator.
func New(input []string) *Generator {
	g := Generator{
		Input: make([]*Input, 0),
	}
	for _, path := range input {
		i := &Input{path: path}
		g.Input = append(g.Input, i)
	}
	return &g
}

// Verbose modifies the verbosity level.
func (g *Generator) Verbose(mode bool) *Generator {
	if mode {
		log.Level = logrus.DebugLevel
	}
	return g
}

// Parse iterates through every input file, calling the proper method.
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
	return nil
}
