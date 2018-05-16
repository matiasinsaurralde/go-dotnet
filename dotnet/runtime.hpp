#pragma once

#ifdef __cplusplus

#include "coreclrhost.h"
#ifndef SUCCEEDED
#define SUCCEEDED(Status) ((Status) >= 0)
#endif // !SUCCEEDED

void* hostHandle;
unsigned int domainId;

void* coreclrLib;
coreclr_initialize_ptr initialize_core_clr;
coreclr_execute_assembly_ptr execute_assembly;
coreclr_shutdown_ptr shutdown_core_clr;
coreclr_create_delegate_ptr create_delegate;

static const char* serverGcVar = "CORECLR_SERVER_GC";
const char* useServerGc;

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
            
int createDelegate(const char* entryPointAssemblyName, const char* entryPointTypeName, const char* entryPointMethodName, int delegateID, void** f);

void parseValues(const char*, char**, int);
#ifdef __cplusplus
}
#endif
