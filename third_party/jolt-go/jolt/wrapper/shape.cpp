/*
 * Jolt Physics C Wrapper - Shape Creation Implementation
 */

#include "shape.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/Collision/Shape/SphereShape.h>
#include <Jolt/Physics/Collision/Shape/BoxShape.h>
#include <Jolt/Physics/Collision/Shape/CapsuleShape.h>
#include <Jolt/Physics/Collision/Shape/ConvexHullShape.h>
#include <Jolt/Physics/Collision/Shape/MeshShape.h>
#include <Jolt/Physics/Collision/RayCast.h>
#include <Jolt/Physics/Collision/CastResult.h>
#include <Jolt/Physics/Collision/Shape/SubShapeID.h>
#include <Jolt/Physics/Collision/TransformedShape.h>

using namespace JPH;

JoltShape JoltCreateSphere(float radius)
{
	SphereShapeSettings sphere_settings(radius);
	ShapeSettings::ShapeResult sphere_result = sphere_settings.Create();

	// Shapes are ref-counted, AddRef to keep it alive
	ShapeRefC shape = sphere_result.Get();
	shape->AddRef();

	return static_cast<JoltShape>(const_cast<Shape*>(shape.GetPtr()));
}

JoltShape JoltCreateBox(float halfExtentX, float halfExtentY, float halfExtentZ)
{
	BoxShapeSettings box_settings(Vec3(halfExtentX, halfExtentY, halfExtentZ));
	ShapeSettings::ShapeResult box_result = box_settings.Create();

	// Shapes are ref-counted, AddRef to keep it alive
	ShapeRefC shape = box_result.Get();
	shape->AddRef();

	return static_cast<JoltShape>(const_cast<Shape*>(shape.GetPtr()));
}

JoltShape JoltCreateCapsule(float halfHeight, float radius)
{
	CapsuleShapeSettings capsule_settings(halfHeight, radius);
	ShapeSettings::ShapeResult capsule_result = capsule_settings.Create();

	// Shapes are ref-counted, AddRef to keep it alive
	ShapeRefC shape = capsule_result.Get();
	shape->AddRef();

	return static_cast<JoltShape>(const_cast<Shape*>(shape.GetPtr()));
}

JoltShape JoltCreateConvexHull(const float* points, int numPoints)
{
	// Convert float array to Vec3 array
	Array<Vec3> vertices;
	vertices.reserve(numPoints);
	for (int i = 0; i < numPoints; ++i) {
		vertices.push_back(Vec3(points[i * 3], points[i * 3 + 1], points[i * 3 + 2]));
	}

	ConvexHullShapeSettings hull_settings(vertices);
	ShapeSettings::ShapeResult hull_result = hull_settings.Create();

	// Shapes are ref-counted, AddRef to keep it alive
	ShapeRefC shape = hull_result.Get();
	shape->AddRef();

	return static_cast<JoltShape>(const_cast<Shape*>(shape.GetPtr()));
}

JoltShape JoltCreateMesh(const float* vertices, int numVertices,
							   const int* indices, int numIndices)
{
	// Create mesh shape with vertices and indices
	TriangleList triangles;
	triangles.reserve(numIndices / 3);

	for (int i = 0; i < numIndices; i += 3) {
		int i0 = indices[i];
		int i1 = indices[i + 1];
		int i2 = indices[i + 2];

		Triangle tri;
		tri.mV[0] = Float3(vertices[i0 * 3], vertices[i0 * 3 + 1], vertices[i0 * 3 + 2]);
		tri.mV[1] = Float3(vertices[i1 * 3], vertices[i1 * 3 + 1], vertices[i1 * 3 + 2]);
		tri.mV[2] = Float3(vertices[i2 * 3], vertices[i2 * 3 + 1], vertices[i2 * 3 + 2]);
		triangles.push_back(tri);
	}

	MeshShapeSettings mesh_settings(triangles);
	ShapeSettings::ShapeResult mesh_result = mesh_settings.Create();

	// Shapes are ref-counted, AddRef to keep it alive
	ShapeRefC shape = mesh_result.Get();
	shape->AddRef();

	return static_cast<JoltShape>(const_cast<Shape*>(shape.GetPtr()));
}

void JoltDestroyShape(JoltShape shape)
{
	Shape* s = static_cast<Shape*>(shape);
	s->Release();
}

int JoltShapeCastRay(JoltShape shape,
                     float originX, float originY, float originZ,
                     float directionX, float directionY, float directionZ,
                     int backfaceMode, int treatConvexAsSolid,
                     float* outFraction)
{
	Shape* s = static_cast<Shape*>(shape);
	if (!s) return 0;

	// Create the ray (convert to single precision RayCast)
	RayCast ray(Vec3(originX, originY, originZ), Vec3(directionX, directionY, directionZ));

	// Create SubShapeIDCreator for tracking hierarchical shape IDs
	SubShapeIDCreator sub_shape_creator;

	// Cast the ray
	RayCastResult result;
	if (s->CastRay(ray, sub_shape_creator, result)) {
		*outFraction = result.mFraction;
		return 1; // Hit
	}

	return 0; // No hit
}

JoltTransformedShape JoltCreateTransformedShape(JoltShape shape,
                                                 float posX, float posY, float posZ,
                                                 float rotX, float rotY, float rotZ, float rotW,
                                                 unsigned int bodyID)
{
	Shape* s = static_cast<Shape*>(shape);
	if (!s) return nullptr;

	// Create a TransformedShape on the heap
	TransformedShape* ts = new TransformedShape(
		RVec3(posX, posY, posZ),
		Quat(rotX, rotY, rotZ, rotW),
		s,
		BodyID(bodyID)
	);

	return static_cast<JoltTransformedShape>(ts);
}

void JoltDestroyTransformedShape(JoltTransformedShape transformedShape)
{
	TransformedShape* ts = static_cast<TransformedShape*>(transformedShape);
	delete ts;
}

int JoltTransformedShapeCastRay(JoltTransformedShape transformedShape,
                                 float originX, float originY, float originZ,
                                 float directionX, float directionY, float directionZ,
                                 float* outFraction)
{
	TransformedShape* ts = static_cast<TransformedShape*>(transformedShape);
	if (!ts) return 0;

	// Create the ray (RRayCast for double precision)
	RRayCast ray(RVec3(originX, originY, originZ), Vec3(directionX, directionY, directionZ));

	// Cast the ray
	RayCastResult result;
	if (ts->CastRay(ray, result)) {
		*outFraction = result.mFraction;
		return 1; // Hit
	}

	return 0; // No hit
}
