package mbgame

func packRGB(r, g, b int64) int64 {
	return ((r & 255) << 16) | ((g & 255) << 8) | (b & 255)
}

func packARGB(a, r, g, b int64) int64 {
	return ((a & 255) << 24) | ((r & 255) << 16) | ((g & 255) << 8) | (b & 255)
}

func unpackR(col int64) int64 { return (col >> 16) & 255 }
func unpackG(col int64) int64 { return (col >> 8) & 255 }
func unpackB(col int64) int64 { return col & 255 }
func unpackA(col int64) int64 { return (col >> 24) & 255 }

func rgbMix(c1, c2 int64, t float64) int64 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	a1, r1, g1, b1 := unpackA(c1), unpackR(c1), unpackG(c1), unpackB(c1)
	a2, r2, g2, b2 := unpackA(c2), unpackR(c2), unpackG(c2), unpackB(c2)
	a := int64(float64(a1) + (float64(a2)-float64(a1))*t)
	r := int64(float64(r1) + (float64(r2)-float64(r1))*t)
	g := int64(float64(g1) + (float64(g2)-float64(g1))*t)
	b := int64(float64(b1) + (float64(b2)-float64(b1))*t)
	return packARGB(a, r, g, b)
}

func rgbBrighten(col int64, amount float64) int64 {
	a, r, g, b := unpackA(col), unpackR(col), unpackG(col), unpackB(col)
	f := 1 + amount
	return packARGB(
		a,
		clampU8(int64(float64(r)*f)),
		clampU8(int64(float64(g)*f)),
		clampU8(int64(float64(b)*f)),
	)
}

func rgbDarken(col int64, amount float64) int64 {
	a, r, g, b := unpackA(col), unpackR(col), unpackG(col), unpackB(col)
	f := 1 - amount
	if f < 0 {
		f = 0
	}
	return packARGB(
		a,
		clampU8(int64(float64(r)*f)),
		clampU8(int64(float64(g)*f)),
		clampU8(int64(float64(b)*f)),
	)
}

func rgbFade(col int64, alpha float64) int64 {
	a := clampU8(int64(alpha * 255))
	_, r, g, b := unpackA(col), unpackR(col), unpackG(col), unpackB(col)
	return packARGB(int64(a), r, g, b)
}

func clampU8(n int64) int64 {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return n
}

func colWhite() int64   { return packRGB(255, 255, 255) }
func colBlack() int64   { return packRGB(0, 0, 0) }
func colRed() int64     { return packRGB(255, 0, 0) }
func colGreen() int64   { return packRGB(0, 255, 0) }
func colBlue() int64    { return packRGB(0, 0, 255) }
func colYellow() int64  { return packRGB(255, 255, 0) }
func colCyan() int64    { return packRGB(0, 255, 255) }
func colMagenta() int64 { return packRGB(255, 0, 255) }
func colOrange() int64  { return packRGB(255, 165, 0) }
func colGray() int64    { return packRGB(128, 128, 128) }
func colDarkGray() int64 { return packRGB(64, 64, 64) }
func colLightGray() int64 { return packRGB(192, 192, 192) }
func colTransparent() int64 { return packARGB(0, 0, 0, 0) }
