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
