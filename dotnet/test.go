package dotnet

/*
typedef int (*AddFunc)(int, int);
AddFunc addFunc;

void** getAddFunc() {
	return (void**)&addFunc;
}

int callAddFunc(int a, int b) {
	return addFunc(a, b);
}

typedef void (*DummyFunc)();
DummyFunc dummyFunc;

void** getDummyFunc() {
	return (void**)&dummyFunc;
}

typedef char* (*StringFunc)();
StringFunc stringFunc;

void** getStringFunc() {
	return (void**)&stringFunc;
}

char* callStringFunc() {
	return stringFunc();
}
*/
import "C"
import "unsafe"

func getAddFunc() *unsafe.Pointer {
	return C.getAddFunc()
}

func callAddFunc(a, b int) int {
	return int(C.callAddFunc(C.int(a), C.int(b)))
}

func getDummyFunc() *unsafe.Pointer {
	return C.getDummyFunc()
}

func getStringFunc() *unsafe.Pointer {
	return C.getStringFunc()
}

func callStringFunc() string {
	return C.GoString(C.callStringFunc())
}
