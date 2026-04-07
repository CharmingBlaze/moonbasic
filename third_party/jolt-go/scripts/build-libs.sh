#!/bin/bash
#
# Build pre-built binaries for all supported platforms
#
# Requirements:
#   - JOLT_SRC environment variable pointing to JoltPhysics source
#   - Docker (for Linux builds)
#   - macOS with ARM64 (for darwin_arm64 builds)
#
# Usage:
#   export JOLT_SRC=/path/to/JoltPhysics
#   ./scripts/build-libs.sh [darwin_arm64|linux_amd64|all]
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}ℹ ${NC}$1"
}

success() {
    echo -e "${GREEN}✅ ${NC}$1"
}

error() {
    echo -e "${RED}❌ ${NC}$1"
}

warn() {
    echo -e "${YELLOW}⚠️  ${NC}$1"
}

# Check JOLT_SRC
if [ -z "$JOLT_SRC" ]; then
    error "JOLT_SRC environment variable is not set"
    echo "Please set it to the path of JoltPhysics source directory:"
    echo "  export JOLT_SRC=/path/to/JoltPhysics"
    exit 1
fi

if [ ! -d "$JOLT_SRC" ]; then
    error "JOLT_SRC directory does not exist: $JOLT_SRC"
    exit 1
fi

info "Using Jolt Physics from: $JOLT_SRC"

# Get Jolt version
JOLT_VERSION=$(cd "$JOLT_SRC" && git describe --tags --always 2>/dev/null || echo "unknown")
info "Jolt version: $JOLT_VERSION"

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
WRAPPER_DIR="$REPO_ROOT/jolt/wrapper"
LIB_DIR="$REPO_ROOT/jolt/lib"

# Determine what to build
TARGET="${1:-all}"

build_darwin_arm64() {
    info "Building darwin/arm64..."

    # Check if we're on macOS ARM64
    if [ "$(uname -s)" != "Darwin" ]; then
        warn "Skipping darwin/arm64 - not on macOS"
        return
    fi

    if [ "$(uname -m)" != "arm64" ]; then
        warn "Skipping darwin/arm64 - not on ARM64 architecture"
        return
    fi

    BUILD_DIR="$JOLT_SRC/Build/macos_arm64_release"

    info "  Building Jolt library..."
    mkdir -p "$BUILD_DIR"
    cd "$BUILD_DIR"
    cmake .. \
        -DCMAKE_BUILD_TYPE=Release \
        -DCMAKE_OSX_ARCHITECTURES=arm64 \
        -DCMAKE_CXX_COMPILER=clang++ \
        -DDISABLE_CUSTOM_ALLOCATOR=ON \
        -DTARGET_UNIT_TESTS=OFF \
        -DTARGET_HELLO_WORLD=OFF \
        -DTARGET_PERFORMANCE_TEST=OFF \
        -DTARGET_SAMPLES=OFF \
        -DTARGET_VIEWER=OFF
    cmake --build . -j$(sysctl -n hw.ncpu)

    info "  Building wrapper..."
    cd "$WRAPPER_DIR"

    # Compile all wrapper source files (auto-discover .cpp files)
    for src in *.cpp; do
        if [ -f "$src" ]; then
            clang++ -std=c++17 \
                -I"$JOLT_SRC" \
                -DNDEBUG \
                -DJPH_DISABLE_CUSTOM_ALLOCATOR \
                -DJPH_PROFILE_ENABLED \
                -DJPH_DEBUG_RENDERER \
                -DJPH_OBJECT_STREAM \
                -c "$src" \
                -o "${src%.cpp}.o"
        fi
    done

    # Create static library from all object files
    ar rcs libjolt_wrapper.a *.o
    rm *.o

    info "  Copying to $LIB_DIR/darwin_arm64..."
    mkdir -p "$LIB_DIR/darwin_arm64"
    cp libjolt_wrapper.a "$LIB_DIR/darwin_arm64/"
    cp "$BUILD_DIR/libJolt.a" "$LIB_DIR/darwin_arm64/"
    rm libjolt_wrapper.a

    success "darwin/arm64 build complete"
    ls -lh "$LIB_DIR/darwin_arm64/"
}

build_linux_amd64() {
    info "Building linux/amd64..."

    # Check if Docker is available
    if ! command -v docker &> /dev/null; then
        error "Docker is required for Linux builds but is not installed"
        exit 1
    fi

    info "  Building Docker image..."
    docker build \
        --platform linux/amd64 \
        -f "$SCRIPT_DIR/docker/Dockerfile.linux-amd64" \
        -t jolt-builder-linux-amd64 \
        "$SCRIPT_DIR/docker"

    info "  Running build in Docker container..."
    mkdir -p "$LIB_DIR/linux_amd64"

    # Copy wrapper files to a temp directory that we can mount
    TEMP_WRAPPER=$(mktemp -d)
    cp "$WRAPPER_DIR"/*.cpp "$WRAPPER_DIR"/*.h "$TEMP_WRAPPER/"

    docker run --rm \
        --platform linux/amd64 \
        -v "$JOLT_SRC:/build/JoltPhysics" \
        -v "$TEMP_WRAPPER:/build/wrapper" \
        -v "$LIB_DIR/linux_amd64:/build/output" \
        jolt-builder-linux-amd64

    # Clean up temp directory
    rm -rf "$TEMP_WRAPPER"

    success "linux/amd64 build complete"
    ls -lh "$LIB_DIR/linux_amd64/"
}

build_linux_arm64() {
    info "Building linux/arm64..."

    # Check if Docker is available
    if ! command -v docker &> /dev/null; then
        error "Docker is required for Linux builds but is not installed"
        exit 1
    fi

    info "  Building Docker image..."
    docker build \
        --platform linux/arm64 \
        -f "$SCRIPT_DIR/docker/Dockerfile.linux-arm64" \
        -t jolt-builder-linux-arm64 \
        "$SCRIPT_DIR/docker"

    info "  Running build in Docker container..."
    mkdir -p "$LIB_DIR/linux_arm64"

    # Copy wrapper files to a temp directory that we can mount
    TEMP_WRAPPER=$(mktemp -d)
    cp "$WRAPPER_DIR"/*.cpp "$WRAPPER_DIR"/*.h "$TEMP_WRAPPER/"

    docker run --rm \
        --platform linux/arm64 \
        -v "$JOLT_SRC:/build/JoltPhysics" \
        -v "$TEMP_WRAPPER:/build/wrapper" \
        -v "$LIB_DIR/linux_arm64:/build/output" \
        jolt-builder-linux-arm64

    # Clean up temp directory
    rm -rf "$TEMP_WRAPPER"

    success "linux/arm64 build complete"
    ls -lh "$LIB_DIR/linux_arm64/"
}

# Main build logic
case "$TARGET" in
    darwin_arm64)
        build_darwin_arm64
        ;;
    linux_amd64)
        build_linux_amd64
        ;;
    linux_arm64)
        build_linux_arm64
        ;;
    all)
        build_darwin_arm64
        build_linux_amd64
        build_linux_arm64
        ;;
    *)
        error "Unknown target: $TARGET"
        echo "Usage: $0 [darwin_arm64|linux_amd64|linux_arm64|all]"
        exit 1
        ;;
esac

echo ""
success "All builds complete! 🎉"
echo ""
info "Next steps:"
echo "  1. Test the binaries: go run example/main.go"
echo "  2. Commit to repo: git add jolt/lib/ && git commit -m 'Update binaries for Jolt $JOLT_VERSION'"
echo ""
