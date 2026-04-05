package mbgame

// Internal state buckets (filled in by helpers; keep types here to avoid circular refs).

type shakeState struct {
	tRem float64
	mag  float64
}

type screenFlashState struct {
	r, g, b, a int
	tRem, dur   float64
}

type musicCrossfadeState struct{}

type burstState struct{}

type vibrateState struct {
	left, right float32
	tRem        float64
}

type fpsCamState struct{}

type tpsCamState struct{}
