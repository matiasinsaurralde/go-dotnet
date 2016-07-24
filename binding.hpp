#ifdef __cplusplus

#include "coreclrhost.h"
#ifndef SUCCEEDED
#define SUCCEEDED(Status) ((Status) >= 0)
#endif // !SUCCEEDED

void* hostHandle;
unsigned int domainId;

extern "C" {
#endif

int initializeCoreCLR(const char* exePath,
            const char* appDomainFriendlyName,
            int propertyCount,
            const char* mergedPropertyKeys,
            const char* mergedPropertyValues,
            const char* managedAssemblyAbsolutePath,
            const char* clrFilesAbsolutePath);
int shutdownCoreCLR();
int executeManagedAssembly(const char*);
int createDelegate(const char* entryPointAssemblyName,
            const char* entryPointTypeName,
            const char* entryPointMethodName, int delegateId);

void parseValues(const char*, char**, int);

#define __stdcall
#define STDMETHODCALLTYPE __stdcall
typedef void (STDMETHODCALLTYPE *TheFunction)();

#ifdef __cplusplus
}
#endif
