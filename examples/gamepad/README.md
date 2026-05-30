# Gamepad example

Demonstrates `INPUT.ISGAMEPADAVAILABLE`, `INPUT.GETGAMEPADAXISVALUE`, and `INPUT.JOYDOWN` with the `GAMEPAD_*` constants seeded at runtime.

```bash
CGO_ENABLED=1 go run . examples/gamepad/main.mb
```

Connect a controller (pad index `0`) before or during the demo. Left stick moves the square; face buttons boost color channels.
