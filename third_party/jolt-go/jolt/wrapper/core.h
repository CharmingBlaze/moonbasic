/*
 * Jolt Physics C Wrapper - Core Initialization
 *
 * Handles global initialization, shutdown, and shared resources.
 */

#ifndef JOLT_WRAPPER_CORE_H
#define JOLT_WRAPPER_CORE_H

#ifdef __cplusplus
extern "C" {
#endif

// Initialize Jolt Physics (call once at startup)
// Returns 1 on success, 0 on failure
int JoltInit();

// Shutdown Jolt Physics (call once at exit)
void JoltShutdown();

#ifdef __cplusplus
}

// C++ only: Access to global resources
#include <memory>

namespace JPH {
    class TempAllocatorImpl;
    class JobSystemThreadPool;
}

extern std::unique_ptr<JPH::TempAllocatorImpl> gTempAllocator;
extern std::unique_ptr<JPH::JobSystemThreadPool> gJobSystem;

#endif

#endif // JOLT_WRAPPER_CORE_H
