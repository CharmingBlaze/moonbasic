package window

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

// envTruthy treats "1", "true", "yes", "on" (case-insensitive) as true.
func envTruthy(key string) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	switch v {
	case "1", "true", "yes", "on":
		return true
	default:
		return false
	}
}

// windowOpenWantHighDPI is true when FLAG_WINDOW_HIGHDPI should be set in SetConfigFlags before InitWindow.
// Default is false so low-end integrated GPUs and Windows display scaling are less likely to show a blank
// or mis-mapped client area. Opt in with MOONBASIC_ENABLE_HIGHDPI=1.
func windowOpenWantHighDPI() bool {
	return envTruthy("MOONBASIC_ENABLE_HIGHDPI")
}

// minimalOpenHandshake skips the post-open presentation guard, message-queue drain, and blank-frame
// warmup (same effect as MOONBASIC_SKIP_OPEN_PRESENT_KICK=1, MOONBASIC_OPEN_WARMUP_FRAMES=0, and no
// extra Poll drain). For diagnosing driver-specific issues on legacy samples; default remains full handshake.
func minimalOpenHandshake() bool {
	return envTruthy("MOONBASIC_MINIMAL_OPEN_HANDSHAKE")
}

// openWarmupBlankFrameCount returns how many full BeginDrawing/clear/EndDrawing cycles to run
// after WINDOW.OPEN (before script code). Set MOONBASIC_OPEN_WARMUP_FRAMES=N (0–120).
// Unset on Windows defaults to 2 (integrated-GPU swap-chain stability); other OS default 0.
// Set MOONBASIC_OPEN_WARMUP_FRAMES=0 to force none on Windows.
func openWarmupBlankFrameCount() int {
	s := strings.TrimSpace(os.Getenv("MOONBASIC_OPEN_WARMUP_FRAMES"))
	if s == "" {
		if runtime.GOOS == "windows" {
			return 2
		}
		return 0
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return 0
	}
	if n > 120 {
		return 120
	}
	return n
}
