package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"html/template"
	"regexp"
	"strconv"
	"strings"
)

const (
	commentPrefix = "// "

	delegatePrefix     = "create_delegate"
	delegateGoTemplate = `
	func {{.MethodName}}({{.Params}}) {{.Returns}} {
		return {{.Call}}
	}
	`

	delegateCTemplate = `
	{{.Returns}} {{.MethodName}}({{.Params}}) {
		return {{.Call}};
	}
	`
)

var (
	delegateExpr = regexp.MustCompile(`(.*)\s(.*)\s(.*)\((.*)\)\s?(.*)`)
)

// DelegateType is used by the code generator to guess equivalent types.
type DelegateType int

// GoType returns the appropriate Golang type.
func (d DelegateType) GoType() (t string) {
	switch d {
	case DelegateIntParam:
		t = "int"
	}
	return t
}

// CType returns the appropriate C type.
func (d DelegateType) CType() (t string) {
	switch d {
	case DelegateIntParam:
		t = "int"
	}
	return t
}

// CGoWrap returns the appropriate CGO wrap.
func (d DelegateType) CGoWrap() (t string) {
	switch d {
	case DelegateIntParam:
		t = "C.int"
	}
	return t
}

// GoWrap returns the appropriate Go wrap.
func (d DelegateType) GoWrap() (t string) {
	switch d {
	case DelegateIntParam:
		t = "int"
	}
	return t
}

const (
	_ DelegateType = iota
	// DelegateIntParam represents the int type.
	DelegateIntParam
)

// DelegateParam contains information about a function param.
type DelegateParam struct {
	Name string
	Type DelegateType
}

// DelegateReturn contains information about a function return value.
type DelegateReturn struct {
	Name string
	Type DelegateType
}

// DelegateAnnotation contains parameters required for generating delegate code
type DelegateAnnotation struct {
	AssemblyName string
	TypeName     string
	MethodName   string

	Params  []DelegateParam
	Returns []DelegateReturn

	*ast.FuncDecl
	*Input
}

type delegateGoTemplateData struct {
	MethodName string
	Params     string
	Returns    string
	Call       string
	Ret        string
}

type delegateCTemplateData struct {
	MethodName string
	Params     string
	Returns    string
	Call       string
	Ret        string
}

type bindingsTemplateData struct {
	HeaderDefinitions string
}

// Render compiles the template using the available info.
// TODO: avoid dup code when matching types
func (d DelegateAnnotation) Render() string {
	// Handle params:
	for _, p := range d.FuncDecl.Type.Params.List {
		paramName := p.Names[0].Name
		param := DelegateParam{
			Name: paramName,
		}
		t := fmt.Sprintf("%v", p.Type)
		switch t {
		case "int":
			param.Type = DelegateIntParam
		default:
			log.Fatalf("Unsupported type '%s' in function '%s'\n", t, d.Name.Name)
		}
		d.Params = append(d.Params, param)
	}

	// Handle returns:
	if d.FuncDecl.Type.Results != nil {
		for _, p := range d.FuncDecl.Type.Results.List {
			var name string
			if len(p.Names) > 0 {
				name = p.Names[0].Name
			}
			result := DelegateReturn{
				Name: name,
			}
			t := fmt.Sprintf("%v", p.Type)
			switch t {
			case "int":
				result.Type = DelegateIntParam
			default:
				log.Fatalf("Unsupported type '%s' in function '%s'\n", t, d.Name.Name)
			}
			d.Returns = append(d.Returns, result)
		}
	}

	// Build cgo call:
	cgoCallName := fmt.Sprintf("createDelegate%s", d.MethodName)
	cgoCallReturns := ""
	for _, v := range d.Returns {
		cgoCallReturns += v.Type.CType()
	}

	cgoCallParams := ""
	for _, v := range d.Params {
		cgoCallParams += v.Type.CType()
	}
	cgoCallDefinition := fmt.Sprintf("%s %s(%s);",
		cgoCallReturns,
		cgoCallName,
		cgoCallParams,
	)

	// typedef int (*HelloWorld)(int);
	typeDefName := fmt.Sprintf("%sFunc", d.MethodName)
	typeDef := fmt.Sprintf("typedef %s (*%s)(%s);", cgoCallReturns, typeDefName, cgoCallParams)

	d.bindingHeaders.WriteString(cgoCallDefinition + "\n")
	d.bindingHeaders.WriteString(typeDef + "\n")

	namedParams := ""
	for _, v := range d.Params {
		namedParams += v.Type.CType() + " " + v.Name
	}
	impl := "\n\t"
	impl += fmt.Sprintf("%s f;", typeDefName)
	impl += "\n\t"
	impl += fmt.Sprintf("create_delegate(hostHandle, domainId, %s, %s, %s, (void**)&f);",
		strconv.Quote(d.AssemblyName),
		strconv.Quote(d.TypeName),
		strconv.Quote(d.MethodName))
	impl += "\n\t"

	letterParams := ""
	for _, v := range d.Params {
		letterParams += v.Name
	}

	impl += fmt.Sprintf("return f(%s);", letterParams)
	impl += "\n"

	cgoCallImpl := fmt.Sprintf("%s %s(%s) {%s}", cgoCallReturns, cgoCallName, namedParams, impl)

	d.bindingSource.WriteString(cgoCallImpl + "\n")

	// Generate cgo wrap cpde:
	out := bytes.Buffer{}
	goTemplateData := delegateGoTemplateData{MethodName: d.MethodName}

	params := ""
	for _, v := range d.Params {
		params += v.Name + " " + v.Type.GoType()
	}
	goTemplateData.Params = params

	returns := ""
	for _, v := range d.Returns {
		returns += v.Type.GoType() + " "
	}
	goTemplateData.Returns = returns

	params = ""
	for _, v := range d.Params {
		params += v.Type.CGoWrap() + "(" + v.Name + ")"
	}

	goTemplateData.Call = fmt.Sprintf("%s(C.%s(%s))", d.Returns[0].Type.GoWrap(), cgoCallName, params)

	cTemplate := template.Must(template.New("delegate_go").Parse(delegateGoTemplate))
	_ = cTemplate.Execute(&out, goTemplateData)

	out.WriteTo(d.goCode)
	return ""
}

// Annotation is an interface.
type Annotation interface {
	Render() string
}

// TODO: implement error handling
func (i *Input) parseFuncAnnotation(s string, f *ast.FuncDecl) Annotation {
	s = strings.TrimLeft(s, commentPrefix)
	var annotation Annotation
	if strings.HasPrefix(s, delegatePrefix) {
		prefix := fmt.Sprintf("//%s:", delegatePrefix)
		s = strings.TrimLeft(s, prefix)
		s = strings.TrimSpace(s)
		submatches := delegateExpr.FindAllStringSubmatch(s, -1)
		matches := submatches[0]
		annotation = &DelegateAnnotation{
			AssemblyName: matches[1],
			TypeName:     matches[2],
			MethodName:   matches[3],
			FuncDecl:     f,
			Input:        i,
		}
	}
	return annotation
}
