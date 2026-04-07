#!/bin/bash
set -e

echo 'Building Jolt Physics for linux/amd64...'

# Build Jolt Physics
mkdir -p /build/JoltPhysics/Build/linux_amd64_release
cd /build/JoltPhysics/Build/linux_amd64_release

cmake .. \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_CXX_COMPILER=g++ \
    -DCMAKE_CXX_FLAGS="-fno-lto" \
    -DCMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF \
    -DDISABLE_CUSTOM_ALLOCATOR=ON \
    -DTARGET_UNIT_TESTS=OFF \
    -DTARGET_HELLO_WORLD=OFF \
    -DTARGET_PERFORMANCE_TEST=OFF \
    -DTARGET_SAMPLES=OFF \
    -DTARGET_VIEWER=OFF

cmake --build . -j$(nproc)

echo 'Building Jolt wrapper for linux/amd64...'

# Build wrapper
cd /build/wrapper

# Compile all wrapper source files (auto-discover .cpp files)
for src in *.cpp; do
    if [ -f "$src" ]; then
        g++ -std=c++17 \
            -I/build/JoltPhysics \
            -DNDEBUG \
            -DJPH_DISABLE_CUSTOM_ALLOCATOR \
            -DJPH_PROFILE_ENABLED \
            -DJPH_DEBUG_RENDERER \
            -DJPH_OBJECT_STREAM \
            -fPIC \
            -fno-lto \
            -c "$src" \
            -o "${src%.cpp}.o"
    fi
done

# Create static library from all object files
ar rcs libjolt_wrapper.a *.o

echo 'Copying binaries to output...'

# Copy binaries
mkdir -p /build/output
cp libjolt_wrapper.a /build/output/
cp /build/JoltPhysics/Build/linux_amd64_release/libJolt.a /build/output/

echo '✅ Linux amd64 binaries built successfully!'
ls -lh /build/output/
