package dotnet

import (
	"path/filepath"
	"runtime"
	"testing"
)

var (
	packagePath  string
	assemblyPath string
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	packagePath = filepath.Dir(filename)
	assemblyPath = filepath.Join(packagePath, "testfiles")
	// copyTestAssemblies()
	SetParams(RuntimeParams{
		Properties: map[string]string{
			"APP_PATHS":                     assemblyPath,
			"NATIVE_DLL_SEARCH_DIRECTORIES": assemblyPath,
		},
	})
	err := Init()
	if err != nil {
		panic(err)
	}
	addFunc := getAddFunc()
	err = CreateDelegate("Test", "Test.TestClass", "Add", 0, addFunc)
	if err != nil {
		panic(err)
	}
	stringFunc := getStringFunc()
	err = CreateDelegate("Test", "Test.TestClass", "String", 0, stringFunc)
	if err != nil {
		panic(err)
	}
}

func TestCreateDelegate(t *testing.T) {
	f := getDummyFunc()
	err := CreateDelegate("foo", "foo.foo", "foo", 1, f)
	if err != errAssemblyNotFound {
		t.Fatalf("Got %s", err.Error())
	}
	err = CreateDelegate("Test", "foo.foo", "foo", 1, f)
	if err != errTypeLoadException {
		t.Fatalf("Got %s", err.Error())
	}
	err = CreateDelegate("Test", "Test.TestClass", "foo", 1, f)
	if err != errMissingMethodException {
		t.Fatalf("Got %s", err.Error())
	}
	err = CreateDelegate("Test", "Test.TestClass", "Add", 1, nil)
	if err != errNullReferenceException {
		t.Fatalf("Got %s", err.Error())
	}
}
func TestAddFunc(t *testing.T) {
	n := callAddFunc(2, 2)
	if n != 4 {
		t.Fatalf("AddFunc call failed, got %d, expected %d", n, 4)
	}
}

func TestStringFunc(t *testing.T) {
	s := callStringFunc()
	if s != "teststring" {
		t.Fatalf("StringFunc call failed, got %s, expected %s", s, "teststring")
	}
}

func BenchmarkAddFunc(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		callAddFunc(2, 2)
	}
}

func BenchmarkStringFunc(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		callStringFunc()
	}
}
