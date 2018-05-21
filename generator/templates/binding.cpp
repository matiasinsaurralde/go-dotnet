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

#include "binding.hpp"

static const char* serverGcVar = "CORECLR_SERVER_GC";
const char* useServerGc;

void* coreclrLib;
coreclr_initialize_ptr initialize_core_clr;
coreclr_execute_assembly_ptr execute_assembly;
coreclr_shutdown_ptr shutdown_core_clr;
coreclr_create_delegate_ptr create_delegate;

{{ .Impls }}
