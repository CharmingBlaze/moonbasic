/*
 * Jolt Physics C Wrapper - Shape Creation
 *
 * Handles creation and destruction of collision shapes.
 */

#ifndef JOLT_WRAPPER_SHAPE_H
#define JOLT_WRAPPER_SHAPE_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types
typedef void* JoltShape;
typedef void* JoltTransformedShape;

// Create a sphere shape
JoltShape JoltCreateSphere(float radius);

// Create a box shape
JoltShape JoltCreateBox(float halfExtentX, float halfExtentY, float halfExtentZ);

// Create a capsule shape
JoltShape JoltCreateCapsule(float halfHeight, float radius);

// Create a convex hull shape from an array of points
JoltShape JoltCreateConvexHull(const float* points, int numPoints);

// Create a mesh shape from vertices and indices
JoltShape JoltCreateMesh(const float* vertices, int numVertices,
                               const int* indices, int numIndices);

// Destroy a shape
void JoltDestroyShape(JoltShape shape);

// Cast a ray against a shape
// Returns 1 if hit, 0 if miss
// outFraction: receives the fraction along the ray where the hit occurred [0, 1]
int JoltShapeCastRay(JoltShape shape,
                     float originX, float originY, float originZ,
                     float directionX, float directionY, float directionZ,
                     int backfaceMode, int treatConvexAsSolid,
                     float* outFraction);

// Create a transformed shape (combines shape with position and rotation)
// bodyID: use 0 if no associated body
JoltTransformedShape JoltCreateTransformedShape(JoltShape shape,
                                                 float posX, float posY, float posZ,
                                                 float rotX, float rotY, float rotZ, float rotW,
                                                 unsigned int bodyID);

// Destroy a transformed shape
void JoltDestroyTransformedShape(JoltTransformedShape transformedShape);

// Cast a ray against a transformed shape in world space
// Returns 1 if hit, 0 if miss
// outFraction: receives the fraction along the ray where the hit occurred [0, 1]
int JoltTransformedShapeCastRay(JoltTransformedShape transformedShape,
                                 float originX, float originY, float originZ,
                                 float directionX, float directionY, float directionZ,
                                 float* outFraction);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_SHAPE_H
