#ifdef __cplusplus
extern "C" {
#endif
void initializeCoreCLR(const char* exePath,
            const char* appDomainFriendlyName,
            int propertyCount,
            const char* propertyKeys,
            const char* propertyValues);
void executeAssembly();
#ifdef __cplusplus
}
#endif
