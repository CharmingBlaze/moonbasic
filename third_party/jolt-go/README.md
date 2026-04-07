# Jolt Physics Go Bindings

Go bindings for [Jolt Physics](https://github.com/jrouwe/JoltPhysics), a fast and flexible physics engine written in C++.

## Features

- Pre-built binaries - no compilation required
- Cross-platform support (macOS ARM64, Linux x86-64, Linux ARM64)
- Simple, idiomatic Go API
- Built on Jolt Physics v5.4.0

## Installation

**Requirements:**
- Go 1.20 or later

```bash
go get github.com/bbitechnologies/jolt-go
```

That's it! Pre-built binaries are included in the repository, so no compilation is required.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/bbitechnologies/jolt-go/jolt"
)

func main() {
    // Initialize Jolt Physics
    if err := jolt.Init(); err != nil {
        panic(err)
    }
    defer jolt.Shutdown()

    // Create physics world
    ps := jolt.NewPhysicsSystem()
    defer ps.Destroy()

    // Create a dynamic sphere
    bi := ps.GetBodyInterface()
    sphere := bi.CreateSphere(
        1.0,                          // radius
        jolt.Vec3{X: 0, Y: 20, Z: 0}, // position
        true,                         // is dynamic
    )
    defer sphere.Destroy()

    // Simulate physics
    for i := 0; i < 60; i++ {
        ps.Update(1.0 / 60.0)

        pos := bi.GetPosition(sphere)
        fmt.Printf("Frame %d: Y = %.2f\n", i, pos.Y)
    }
}
```

See the [example](example/main.go) for a complete working demo.

## Supported Platforms

| Platform      | Architecture | Status |
|---------------|--------------|--------|
| macOS         | ARM64        | ✅     |
| Linux         | x86-64       | ✅     |
| Linux         | ARM64        | ✅     |
| macOS         | x86-64       | 🚧     |
| Windows       | x86-64       | 🚧     |

Want support for another platform? Open an issue or see [CONTRIBUTORS.md](CONTRIBUTORS.md) for build instructions.

## Architecture

```
Go Application
     ↓
CGo Bindings (jolt/*.go)
     ↓
C Wrapper (jolt/wrapper/*.{cpp,h})
     ↓
Jolt Physics C++ Library
```

The wrapper uses opaque pointers to cleanly separate Go and C++ code.

## Performance

Jolt Physics is one of the fastest physics engines available:
- Optimized for modern CPUs with SIMD
- Efficient broad-phase collision detection
- Continuous collision detection (CCD)
- Stable and deterministic simulation

## Contributing

Contributions are welcome! Please see [CONTRIBUTORS.md](CONTRIBUTORS.md) for:
- How to rebuild binaries
- Adding support for new platforms
- Modifying the C wrapper

## Resources

- [Jolt Physics Documentation](https://jrouwe.github.io/JoltPhysics/)
- [Jolt Physics GitHub](https://github.com/jrouwe/JoltPhysics)
- [Jolt Physics Samples](https://github.com/jrouwe/JoltPhysics/tree/master/Samples)

## License

This project follows the same license as Jolt Physics. See [LICENSE](LICENSE) for details.

## Acknowledgments

Built on top of the excellent [Jolt Physics](https://github.com/jrouwe/JoltPhysics) engine by Jorrit Rouwe.
