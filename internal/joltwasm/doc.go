// Package joltwasm hosts the experimental WebAssembly path for 3D physics (Jolt compiled to wasm,
// driven by [github.com/tetratelabs/wazero]). Nothing here is wired into moonBASIC’s runtime yet;
// it exists to validate load/instantiate cost and linear-memory read patterns ahead of a real Jolt build.
//
// Tests embed a small WASI module from the wazero tree (see testdata/) to exercise instantiation;
// replace with a Jolt wasm once the compile pipeline exists. [StateView] documents the SoA read path.
//
// [github.com/tetratelabs/wazero]: https://github.com/tetratelabs/wazero
package joltwasm
