#ifdef __cplusplus

#include "coreclrhost.h"
#ifndef SUCCEEDED
#define SUCCEEDED(Status) ((Status) >= 0)
#endif // !SUCCEEDED

void* hostHandle;
unsigned int domainId;

extern "C" {
#endif

{{ .HeaderDefinitions }}

#ifdef __cplusplus
}
#endif
