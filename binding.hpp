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
void parseValues(const char*, char**, int);
void executeAssembly();
#ifdef __cplusplus
}
#endif
