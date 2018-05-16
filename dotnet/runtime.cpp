#include <stdio.h>
#include <cstdlib>
#include <sstream>
#include <dlfcn.h>
#include <limits.h>
#include <string>
#include <cstring>

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include "coreruncommon.h"

#include "runtime.hpp"

int initializeCoreCLR(const char* exePath,
            const char* appDomainFriendlyName,
            int propertyCount,
            const char* mergedPropertyKeys,
            const char* mergedPropertyValues,
            const char* managedAssemblyAbsolutePath,
            const char* clrFilesAbsolutePath) {

  std::string coreClrDllPath(clrFilesAbsolutePath);
  coreClrDllPath.append("/");
  coreClrDllPath.append(coreClrDll);

  if (coreClrDllPath.length() >= PATH_MAX)
  {
      fprintf(stderr, "Absolute path to libcoreclr.so too long\n");
  }

  std::string appPath;

  if( managedAssemblyAbsolutePath[0] == '\0' ) {
    // printf("Expecting to run a standard .exe\n");
  } else {
    // printf("Expecting to load an assembly and invoke arbitrary methods.\n");
    GetDirectory(managedAssemblyAbsolutePath, appPath);
  };

  // Construct native search directory paths
  std::string nativeDllSearchDirs(appPath);
  char *coreLibraries = getenv("CORE_LIBRARIES");
  if (coreLibraries)
  {
      nativeDllSearchDirs.append(":");
      nativeDllSearchDirs.append(coreLibraries);
  }
  nativeDllSearchDirs.append(":");
  nativeDllSearchDirs.append(clrFilesAbsolutePath);

  std::string tpaList;
  AddFilesFromDirectoryToTpaList(clrFilesAbsolutePath, tpaList);

  coreclrLib = dlopen(coreClrDllPath.c_str(), RTLD_NOW | RTLD_LOCAL);
  if (coreclrLib != nullptr)
  {
      initialize_core_clr = (coreclr_initialize_ptr)dlsym(coreclrLib, "coreclr_initialize");
      execute_assembly = (coreclr_execute_assembly_ptr)dlsym(coreclrLib, "coreclr_execute_assembly");
      shutdown_core_clr= (coreclr_shutdown_ptr)dlsym(coreclrLib, "coreclr_shutdown");
      create_delegate = (coreclr_create_delegate_ptr)dlsym(coreclrLib, "coreclr_create_delegate");

      if (initialize_core_clr == nullptr)
      {
          fprintf(stderr, "Function coreclr_initialize not found in the libcoreclr.so\n");
          return -1;
      }
      else if (execute_assembly == nullptr)
      {
          fprintf(stderr, "Function coreclr_execute_assembly not found in the libcoreclr.so\n");
          return -1;
      }
      else if (shutdown_core_clr == nullptr)
      {
          fprintf(stderr, "Function coreclr_shutdown not found in the libcoreclr.so\n");
          return -1;
      } else {
        if(useServerGc == NULL) {
          std::getenv(serverGcVar);
          if (useServerGc == nullptr) {
              useServerGc = "0";
          }
        }

        useServerGc = std::strcmp(useServerGc, "1") == 0 ? "true" : "false";

        // Keep enough space for inserting the tpaList:
        propertyCount++;

        char *keys[propertyCount];
        char *values[propertyCount];

        parseValues(mergedPropertyKeys, keys, propertyCount);
        parseValues(mergedPropertyValues, values, propertyCount);

        bool tpaOverride = false;

        const char *tpaKey = "TRUSTED_PLATFORM_ASSEMBLIES";

        for( int i = 0; i < propertyCount ; i++ ) {
          int match = strncmp( tpaKey, keys[i], strlen(tpaKey) );
          if( match == 0 ) {
            tpaOverride = true;
            break;
          }
        };

        if( !tpaOverride ) {
          keys[propertyCount] = (char*)std::malloc(strlen(tpaKey)+1);
          std::strcpy(keys[propertyCount], tpaKey);

          values[propertyCount] = (char*)std::malloc(strlen(tpaList.c_str())+1);
          std::strcpy(values[propertyCount], tpaList.c_str());

          propertyCount++;
        };

        int st = initialize_core_clr(
                    exePath,
                    appDomainFriendlyName,
                    propertyCount,
                    (const char**)keys,
                    (const char**)values,
                    &hostHandle,
                    &domainId);

        if (SUCCEEDED(st)) {
          // printf("coreclr_initialize ok\n");
        } else {
          fprintf(stderr, "coreclr_initialize failed - status: 0x%08x\n", st);
        };

      }
    }

    return 0;

}

int shutdownCoreCLR() {
  int st = shutdown_core_clr(hostHandle, domainId);
  if (!SUCCEEDED(st)) {
    fprintf(stderr, "coreclr_shutdown failed - status: 0x%08x\n", st);
    return -1;
  }
  return st;
};

int executeManagedAssembly(const char *assembly) {
  printf("Executing: %s\n", assembly);

  unsigned int* exitCode;
  int st = execute_assembly(
          hostHandle,
          domainId,
          0,
          NULL,
          assembly,
          (unsigned int*)&exitCode);

  if (!SUCCEEDED(st)) {
    return st;
  };
  return 0;
};

void parseValues(const char* input, char** dest, int count) {
  std::stringstream values(input);
  std::string e;

  int i = 0;
  while( std::getline(values, e, ';')) {
    const char *v = e.c_str();
    dest[i] = (char*)std::malloc(strlen(v)+1);
    std::strcpy(dest[i], v);
    i++;
  }
};

 int createDelegate(const char* entryPointAssemblyName, const char* entryPointTypeName, const char* entryPointMethodName, int delegateID, void** f) {
   return create_delegate(hostHandle, domainId, entryPointAssemblyName, entryPointTypeName, entryPointMethodName, f);
 }