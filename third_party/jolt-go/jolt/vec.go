package jolt

import "math"

// Vec3 represents a 3D vector
type Vec3 struct {
	X, Y, Z float32
}

// Add returns the sum of this vector and another vector
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{X: v.X + other.X, Y: v.Y + other.Y, Z: v.Z + other.Z}
}

// Sub returns the difference of this vector and another vector
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{X: v.X - other.X, Y: v.Y - other.Y, Z: v.Z - other.Z}
}

// Mul returns this vector multiplied by a scalar
func (v Vec3) Mul(scalar float32) Vec3 {
	return Vec3{X: v.X * scalar, Y: v.Y * scalar, Z: v.Z * scalar}
}

// Dot returns the dot product of this vector with another vector
func (v Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Length returns the magnitude (length) of the vector
func (v Vec3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Normalize returns a unit vector in the same direction as this vector.
// Returns a zero vector if the input vector has zero length.
func (v Vec3) Normalize() Vec3 {
	length := v.Length()
	if length == 0 {
		return Vec3{X: 0, Y: 0, Z: 0}
	}
	return Vec3{X: v.X / length, Y: v.Y / length, Z: v.Z / length}
}

// Quat represents a quaternion for rotations
type Quat struct {
	X, Y, Z, W float32
}

// Identity returns an identity quaternion (no rotation)
func QuatIdentity() Quat {
	return Quat{X: 0, Y: 0, Z: 0, W: 1}
}
