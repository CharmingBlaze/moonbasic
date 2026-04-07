package jolt

import (
	"math"
	"testing"
)

func TestShapeCastRay(t *testing.T) {
	// Create a sphere shape at origin with radius 1.0
	sphere := CreateSphere(1.0)
	defer sphere.Destroy()

	settings := DefaultRayCastSettings()

	t.Run("Ray hits sphere from above", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 10, Z: 0},
			Direction: Vec3{X: 0, Y: -20, Z: 0}, // Ray goes down 20 units
		}
		result := RayCastResult{}

		if !sphere.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the sphere")
		}

		hitDistance := result.Fraction * 20.0
		expectedDistance := float32(9.0) // 10 (start) - 1 (radius) = 9
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})

	t.Run("Ray misses sphere", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 5, Y: 0, Z: 0},
			Direction: Vec3{X: 0, Y: 10, Z: 0}, // Ray goes up, sphere is at origin
		}
		result := RayCastResult{}

		if sphere.CastRay(ray, settings, &result) {
			t.Error("Ray should have missed the sphere")
		}
	})

	t.Run("Ray hits sphere from side", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: -5, Y: 0, Z: 0},
			Direction: Vec3{X: 10, Y: 0, Z: 0}, // Ray goes right 10 units
		}
		result := RayCastResult{}

		if !sphere.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the sphere")
		}

		hitDistance := result.Fraction * 10.0
		expectedDistance := float32(4.0) // 5 (start) - 1 (radius) = 4
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})

	t.Run("Ray starts inside sphere", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 0, Z: 0}, // Inside the sphere
			Direction: Vec3{X: 10, Y: 0, Z: 0},
		}
		result := RayCastResult{}

		// With TreatConvexAsSolid = true, ray starting inside should hit at fraction 0
		if !sphere.CastRay(ray, settings, &result) {
			t.Fatal("Ray starting inside should hit the sphere")
		}

		if result.Fraction > 0.01 {
			t.Errorf("Ray starting inside should hit at fraction ~0, got %.4f", result.Fraction)
		}
	})
}

func TestShapeCastRayBox(t *testing.T) {
	// Create a box shape centered at origin with half-extents (1, 2, 3)
	box := CreateBox(Vec3{X: 1, Y: 2, Z: 3})
	defer box.Destroy()

	settings := DefaultRayCastSettings()

	t.Run("Ray hits box from above", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 10, Z: 0},
			Direction: Vec3{X: 0, Y: -20, Z: 0},
		}
		result := RayCastResult{}

		if !box.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the box")
		}

		hitDistance := result.Fraction * 20.0
		expectedDistance := float32(8.0) // 10 (start) - 2 (half-extent Y) = 8
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})
}

func TestShapeCastRayCapsule(t *testing.T) {
	// Create a capsule shape centered at origin
	capsule := CreateCapsule(2.0, 0.5) // halfHeight=2.0, radius=0.5
	defer capsule.Destroy()

	settings := DefaultRayCastSettings()

	t.Run("Ray hits capsule", func(t *testing.T) {
		ray := RRayCast{
			Origin:    Vec3{X: -5, Y: 0, Z: 0},
			Direction: Vec3{X: 10, Y: 0, Z: 0},
		}
		result := RayCastResult{}

		if !capsule.CastRay(ray, settings, &result) {
			t.Fatal("Ray should have hit the capsule")
		}

		hitDistance := result.Fraction * 10.0
		// Capsule extends 0.5 units in X direction (radius), so hit at 5 - 0.5 = 4.5
		if hitDistance < 4.0 || hitDistance > 5.0 {
			t.Errorf("Hit distance = %.2f, expected between 4.0 and 5.0", hitDistance)
		}
	})
}

func TestTransformedShapeCastRay(t *testing.T) {
	// Create a sphere shape with radius 1.0
	sphere := CreateSphere(1.0)
	defer sphere.Destroy()

	t.Run("Ray hits transformed shape in world space", func(t *testing.T) {
		// Create a transformed shape positioned at (0, 5, 0) with no rotation
		pos := Vec3{X: 0, Y: 5, Z: 0}
		rot := QuatIdentity()
		transformedShape := CreateTransformedShape(sphere, pos, rot, 0)
		defer transformedShape.Destroy()

		// Cast ray from above in world space
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 10, Z: 0},
			Direction: Vec3{X: 0, Y: -20, Z: 0}, // Ray goes down 20 units
		}
		result := RayCastResult{}

		if !transformedShape.CastRay(ray, &result) {
			t.Fatal("Ray should have hit the transformed shape")
		}

		hitDistance := result.Fraction * 20.0
		// Shape is at Y=5 with radius 1, so top is at Y=6
		// Ray starts at Y=10, so distance is 10 - 6 = 4
		expectedDistance := float32(4.0)
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})

	t.Run("Ray misses transformed shape", func(t *testing.T) {
		// Create a transformed shape positioned at (10, 0, 0) with no rotation
		pos := Vec3{X: 10, Y: 0, Z: 0}
		rot := QuatIdentity()
		transformedShape := CreateTransformedShape(sphere, pos, rot, 0)
		defer transformedShape.Destroy()

		// Cast ray that misses the shape
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 10, Z: 0},
			Direction: Vec3{X: 0, Y: -20, Z: 0}, // Ray goes down at origin, shape is at X=10
		}
		result := RayCastResult{}

		if transformedShape.CastRay(ray, &result) {
			t.Error("Ray should have missed the transformed shape")
		}
	})

	t.Run("Ray hits transformed shape from side", func(t *testing.T) {
		// Create a transformed shape positioned at (5, 0, 0) with no rotation
		pos := Vec3{X: 5, Y: 0, Z: 0}
		rot := QuatIdentity()
		transformedShape := CreateTransformedShape(sphere, pos, rot, 0)
		defer transformedShape.Destroy()

		// Cast ray from left in world space
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 0, Z: 0},
			Direction: Vec3{X: 10, Y: 0, Z: 0}, // Ray goes right 10 units
		}
		result := RayCastResult{}

		if !transformedShape.CastRay(ray, &result) {
			t.Fatal("Ray should have hit the transformed shape")
		}

		hitDistance := result.Fraction * 10.0
		// Shape is at X=5 with radius 1, so left edge is at X=4
		// Ray starts at X=0, so distance is 4
		expectedDistance := float32(4.0)
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})

	t.Run("Multiple transformed shapes", func(t *testing.T) {
		// Create two transformed shapes at different positions
		pos1 := Vec3{X: 0, Y: 5, Z: 0}
		pos2 := Vec3{X: 0, Y: 15, Z: 0}
		rot := QuatIdentity()

		ts1 := CreateTransformedShape(sphere, pos1, rot, 0)
		defer ts1.Destroy()

		ts2 := CreateTransformedShape(sphere, pos2, rot, 0)
		defer ts2.Destroy()

		// Ray from above should hit ts2 first
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 20, Z: 0},
			Direction: Vec3{X: 0, Y: -30, Z: 0},
		}

		result1 := RayCastResult{}
		result2 := RayCastResult{}

		hit1 := ts1.CastRay(ray, &result1)
		hit2 := ts2.CastRay(ray, &result2)

		if !hit1 || !hit2 {
			t.Fatal("Ray should hit both transformed shapes")
		}

		// ts2 is closer, should have smaller fraction
		if result2.Fraction >= result1.Fraction {
			t.Errorf("ts2 should be hit first (fraction %.2f) before ts1 (fraction %.2f)",
				result2.Fraction, result1.Fraction)
		}
	})
}

func TestTransformedShapeWithBox(t *testing.T) {
	// Create a box shape with half-extents (1, 2, 3)
	box := CreateBox(Vec3{X: 1, Y: 2, Z: 3})
	defer box.Destroy()

	t.Run("Ray hits transformed box", func(t *testing.T) {
		// Position the box at (0, 10, 0)
		pos := Vec3{X: 0, Y: 10, Z: 0}
		rot := QuatIdentity()
		transformedBox := CreateTransformedShape(box, pos, rot, 0)
		defer transformedBox.Destroy()

		// Cast ray from above
		ray := RRayCast{
			Origin:    Vec3{X: 0, Y: 20, Z: 0},
			Direction: Vec3{X: 0, Y: -30, Z: 0},
		}
		result := RayCastResult{}

		if !transformedBox.CastRay(ray, &result) {
			t.Fatal("Ray should have hit the transformed box")
		}

		hitDistance := result.Fraction * 30.0
		// Box is at Y=10 with half-extent Y=2, so top is at Y=12
		// Ray starts at Y=20, so distance is 20 - 12 = 8
		expectedDistance := float32(8.0)
		if math.Abs(float64(hitDistance-expectedDistance)) > 0.1 {
			t.Errorf("Hit distance = %.2f, expected ~%.2f", hitDistance, expectedDistance)
		}
	})
}
