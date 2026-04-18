# Build Jolt Physics for Windows x64 (static libs)
# 
# Prerequisites:
# - CMake
# - A C++ compiler toolchain (e.g. MinGW-w64 or MSVC)
# - Jolt Physics source (set JPH_SRC env var)
#
# Usage:
#   $env:JPH_SRC = "C:\path\to\JoltPhysics"
#   .\scripts\build-libs-windows.ps1

$ErrorActionPreference = "Stop"

$REPO_ROOT = Get-Item $PSScriptRoot\..
$WRAPPER_SRC = Join-Path $REPO_ROOT "jolt\wrapper"
$LIB_OUT = Join-Path $REPO_ROOT "jolt\lib\windows_amd64"
$JPH_SRC = $env:JPH_SRC

if (-not $JPH_SRC) {
    Write-Error "JPH_SRC environment variable not set. Please set it to the Jolt Physics source directory."
}

if (-not (Test-Path $JPH_SRC)) {
    Write-Error "Jolt Physics source directory not found: $JPH_SRC"
}

Write-Host "--- Building Jolt Core ---" -ForegroundColor Blue
$BUILD_DIR = Join-Path $JPH_SRC "Build\windows_amd64_release"
if (-not (Test-Path $BUILD_DIR)) { New-Item -ItemType Directory $BUILD_DIR }

Push-Location $BUILD_DIR
# Jolt v5.x: CMakeLists.txt lives under Build/ — configure from Build/<profile> with source ..
# -fno-lto: keep libJolt.a as plain objects so any MinGW gcc can link it (avoids
# "bytecode stream ... LTO version X instead of Y" when the archive was built with another GCC).
cmake .. `
    -DCMAKE_BUILD_TYPE=Release `
    -DCMAKE_CXX_FLAGS="-fno-lto" `
    -DCMAKE_C_FLAGS="-fno-lto" `
    -DDISABLE_CUSTOM_ALLOCATOR=ON `
    -DTARGET_UNIT_TESTS=OFF `
    -DTARGET_HELLO_WORLD=OFF `
    -DTARGET_PERFORMANCE_TEST=OFF `
    -DTARGET_SAMPLES=OFF `
    -DTARGET_VIEWER=OFF
cmake --build . --config Release -j $env:NUMBER_OF_PROCESSORS
Pop-Location

Write-Host "--- Building Wrapper ---" -ForegroundColor Blue
if (-not (Test-Path $LIB_OUT)) { New-Item -ItemType Directory $LIB_OUT -Force }

# Note: Using MinGW-style compilation. For MSVC, one might need a different approach.
# We assume MinGW (g++) is on the path if CGO is enabled on Windows.
$OBJS = @()
$CPPS = Get-ChildItem "$WRAPPER_SRC\*.cpp"
foreach ($src in $CPPS) {
    $obj = $src.FullName.Replace(".cpp", ".o")
    Write-Host "Compiling $($src.Name)..."
    & g++ -std=c++17 -O3 -fno-lto -I"$JPH_SRC" -DNDEBUG -DJPH_DISABLE_CUSTOM_ALLOCATOR -DJPH_PROFILE_ENABLED -DJPH_DEBUG_RENDERER -DJPH_OBJECT_STREAM -c "$($src.FullName)" -o "$obj"
    $OBJS += $obj
}

& ar rcs (Join-Path $LIB_OUT "libjolt_wrapper.a") $OBJS
Remove-Item $OBJS

# Copy Jolt Core
$JOLT_LIB = Get-ChildItem -Path $BUILD_DIR -Filter "libJolt.a" -Recurse | Select-Object -First 1
if (-not $JOLT_LIB) {
    # Try .lib for MSVC
    $JOLT_LIB = Get-ChildItem -Path $BUILD_DIR -Filter "Jolt.lib" -Recurse | Select-Object -First 1
}

if ($JOLT_LIB) {
    Copy-Item $JOLT_LIB.FullName $LIB_OUT
    Write-Host "Success! Binaries placed in $LIB_OUT" -ForegroundColor Green
} else {
    Write-Error "Could not find Jolt static library in build directory."
}
