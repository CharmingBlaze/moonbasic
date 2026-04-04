# Building moonBASIC from Source

This guide provides detailed instructions for compiling the `moonBASIC` interpreter from its source code.

---

## Dependencies

Before you can build, you need the following software installed on your system.

### All Systems
- **Go**: Version 1.22 or later. You can download it from the [official Go website](https://go.dev/dl/).
- **Git**: For cloning the repository.

### Windows
- **A C Compiler**: `moonBASIC` relies on `raylib` and other C libraries, so a C compiler is required. We recommend **MinGW-w64**.
  1.  Install **MSYS2** from [https://www.msys2.org/](https://www.msys2.org/).
  2.  Open the MSYS2 MINGW64 terminal and install the GCC toolchain:
      ```bash
      pacman -S mingw-w64-x86_64-gcc
      ```

### Linux (Debian / Ubuntu)
- **A C Compiler and Libraries**: You'll need `gcc` and the development headers for the libraries `raylib` depends on.
  ```bash
  sudo apt-get update
  sudo apt-get install -y gcc libgl1-mesa-dev libxi-dev \
    libxcursor-dev libxrandr-dev libxinerama-dev \
    libwayland-dev libxkbcommon-dev
  ```

---

## Build Steps

### 1. Clone the Repository

First, get the source code from GitHub:
```bash
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic
```

### 2. Build on Windows

Open a standard Command Prompt (`cmd.exe`) or PowerShell. You must tell Go where to find the C compiler you installed.

```bat
REM Set the CGO_ENABLED flag to allow Go to call C code
set CGO_ENABLED=1

REM Point to the MinGW GCC compiler (adjust path if you installed MSYS2 elsewhere)
set CC=C:\msys64\mingw64\bin\gcc.exe

REM Build the executable
go build -o moonbasic.exe .
```

### 3. Build on Linux

Open a terminal and run the following commands:

```bash
# Set the CGO_ENABLED flag
export CGO_ENABLED=1

# Build the executable
go build -o moonbasic .
```

After a successful build, you can run the interpreter directly or add it to your system's PATH to run it from any directory.
