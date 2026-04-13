package hal

// SystemDevice abstracts windowing and engine life cycle.
type SystemDevice interface {
	InitWindow(width, height int, title string)
	CloseWindow()
	WindowShouldClose() bool
	SetTargetFPS(fps int)
	GetFPS() int
	GetFrameTime() float32
	PollInputEvents()
	SetWindowSize(width, height int)
	GetScreenWidth() int
	GetScreenHeight() int
}
