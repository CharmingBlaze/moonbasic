# Asset Lifecycle & Memory Management

The memory safety strategy ensures custom pointers never natively interact causing dangling C-references inside the VM logic. All assets follow the Handle-Based Architecture.

## 1. Load Requests
Invoking `Scene.Load` or `Entity.Load` triggers parsing sequences through the purely-Go abstract layers (`qmuntal/gltf`), producing Go objects structurally defining layouts before interacting with external hardware APIs.

## 2. Dynamic Handles (Registry[int32]Asset)
No pointer directly propagates towards the `moonBASIC` loop structures. Instead, handles mapping exactly towards the Host internal Go `Registry` are presented. 

Whenever the script executes `Entity.Draw(handle)`:
1. Handlers inside Go `mbentity/mbmodel3d` modules securely authenticate the Handle against the Registry index.
2. If the Handle matches an existent bounds array, operations proceed. 
3. If out-of-bounds, the system utilizes "Hard Guards" clamping parameters, falling straight into failure handling errors rather than panicking against nil-segments!

## 3. Destruction Lifecycles
When the program signals `Shader.Free()` or VM termination invokes the GC lifecycle hook:
1. Target handles evaluate true referencing inside the `Registry`.
2. Safe routines trigger C-bound hooks dropping memory contexts gracefully explicitly checking if dynamic handles were securely cached.
3. Handle ID increments, freeing previous addresses preventing "Double-Free Use-After" bugs identically inside the Engine VM structure loops.
