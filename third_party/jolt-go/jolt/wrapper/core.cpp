/*
 * Jolt Physics C Wrapper - Core Initialization Implementation
 */

#include "core.h"
#include <Jolt/Jolt.h>
#include <Jolt/RegisterTypes.h>
#include <Jolt/Core/Factory.h>
#include <Jolt/Core/TempAllocator.h>
#include <Jolt/Core/JobSystemThreadPool.h>
#include <Jolt/Physics/PhysicsSettings.h>
#include <iostream>
#include <memory>
#include <cstdarg>

using namespace JPH;

// Trace callback for Jolt error messages (helps with debugging)
static void TraceImpl(const char *inFMT, ...)
{
	va_list list;
	va_start(list, inFMT);
	char buffer[1024];
	vsnprintf(buffer, sizeof(buffer), inFMT, list);
	va_end(list);
	std::cout << buffer << std::endl;
}

#ifdef JPH_ENABLE_ASSERTS
static bool AssertFailedImpl(const char *inExpression, const char *inMessage, const char *inFile, uint inLine)
{
	std::cout << inFile << ":" << inLine << ": (" << inExpression << ") "
			  << (inMessage != nullptr ? inMessage : "") << std::endl;
	return true;
};
#endif

// Global Jolt resources (shared by all PhysicsSystems)
// Using smart pointers for automatic cleanup and exception safety
std::unique_ptr<TempAllocatorImpl> gTempAllocator;
std::unique_ptr<JobSystemThreadPool> gJobSystem;
static std::unique_ptr<Factory> gFactory;

int JoltInit()
{
	Trace = TraceImpl;
	JPH_IF_ENABLE_ASSERTS(AssertFailed = AssertFailedImpl;)

	gFactory = std::make_unique<Factory>();
	Factory::sInstance = gFactory.get();
	RegisterTypes();

	gTempAllocator = std::make_unique<TempAllocatorImpl>(10 * 1024 * 1024);
	gJobSystem = std::make_unique<JobSystemThreadPool>(cMaxPhysicsJobs, cMaxPhysicsBarriers,
													   std::thread::hardware_concurrency() - 1);

	return 1;
}

void JoltShutdown()
{
	gJobSystem.reset();
	gTempAllocator.reset();
	gFactory.reset();
	Factory::sInstance = nullptr;
}
