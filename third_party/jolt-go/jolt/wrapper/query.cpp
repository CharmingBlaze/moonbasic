/*
 * Jolt Physics C Wrapper - Collision Query Implementation
 */

#include "query.h"
#include "physics.h"
#include <Jolt/Jolt.h>
#include <Jolt/Physics/PhysicsSystem.h>
#include <Jolt/Physics/Collision/CollideShape.h>
#include <Jolt/Physics/Collision/CollisionCollectorImpl.h>
#include <Jolt/Physics/Collision/NarrowPhaseQuery.h>
#include <Jolt/Physics/Collision/RayCast.h>
#include <Jolt/Physics/Collision/CastResult.h>
#include <Jolt/Physics/Body/BodyID.h>
#include <Jolt/Physics/Body/BodyLockInterface.h>
#include <vector>
#include <algorithm>

using namespace JPH;

// Collision layers (defined in physics.cpp)
namespace Layers
{
	static constexpr ObjectLayer MOVING = 1;
};

// Adapter: converts ObjectVsBroadPhaseLayerFilter to BroadPhaseLayerFilter for queries
class BroadPhaseLayerFilterAdapter : public BroadPhaseLayerFilter
{
public:
	BroadPhaseLayerFilterAdapter(const ObjectVsBroadPhaseLayerFilter* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(BroadPhaseLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectVsBroadPhaseLayerFilter* m_filter;
	ObjectLayer m_object_layer;
};

// Adapter: converts ObjectLayerPairFilter to ObjectLayerFilter for queries
class ObjectLayerFilterAdapter : public ObjectLayerFilter
{
public:
	ObjectLayerFilterAdapter(const ObjectLayerPairFilter* filter, ObjectLayer layer)
		: m_filter(filter), m_object_layer(layer) {}

	virtual bool ShouldCollide(ObjectLayer inLayer) const override
	{
		return m_filter->ShouldCollide(m_object_layer, inLayer);
	}

private:
	const ObjectLayerPairFilter* m_filter;
	ObjectLayer m_object_layer;
};

// Collector that just checks if any collision occurred
class AnyHitCollector : public CollideShapeCollector
{
public:
	AnyHitCollector() : m_hasHit(false) {}

	virtual void AddHit(const CollideShapeResult& inResult) override
	{
		m_hasHit = true;
	}

	bool HasHit() const { return m_hasHit; }

private:
	bool m_hasHit;
};

// Collector that stores all collision hits
class AllHitsCollector : public CollideShapeCollector
{
public:
	AllHitsCollector(JoltCollisionHit* outHits, int maxHits)
		: m_outHits(outHits), m_maxHits(maxHits), m_numHits(0) {}

	virtual void AddHit(const CollideShapeResult& inResult) override
	{
		if (m_numHits < m_maxHits)
		{
			JoltCollisionHit& hit = m_outHits[m_numHits];

			// Store body ID
			BodyID* bodyIDCopy = new BodyID(inResult.mBodyID2);
			hit.bodyID = static_cast<JoltBodyID>(bodyIDCopy);

			// Store contact point (using contact point on second shape)
			Vec3 contactPoint = inResult.mContactPointOn2;
			hit.contactPointX = contactPoint.GetX();
			hit.contactPointY = contactPoint.GetY();
			hit.contactPointZ = contactPoint.GetZ();

			// Store penetration depth
			hit.penetrationDepth = inResult.mPenetrationDepth;

			m_numHits++;
		}
	}

	int GetNumHits() const { return m_numHits; }

private:
	JoltCollisionHit* m_outHits;
	int m_maxHits;
	int m_numHits;
};

int JoltCollideShape(JoltPhysicsSystem system, JoltShape shape,
                     float posX, float posY, float posZ, float penetrationTolerance)
{
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	const Shape* s = static_cast<const Shape*>(shape);

	// Get narrow phase query interface
	const NarrowPhaseQuery& query = ps->GetNarrowPhaseQuery();

	// Create filter adapters (query shape acts as MOVING layer)
	BroadPhaseLayerFilterAdapter bpFilter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter objFilter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Create collector to check for any hit
	AnyHitCollector collector;

	// Setup collision settings
	CollideShapeSettings settings;
	settings.mPenetrationTolerance = penetrationTolerance;

	// Perform collision query
	query.CollideShape(
		s,
		Vec3::sReplicate(1.0f),  // Scale
		RMat44::sTranslation(RVec3(posX, posY, posZ)),  // Transform (position, no rotation)
		settings,
		RVec3::sZero(),  // Base offset
		collector,
		bpFilter,
		objFilter
	);

	return collector.HasHit() ? 1 : 0;
}

int JoltCollideShapeGetHits(JoltPhysicsSystem system, JoltShape shape,
                            float posX, float posY, float posZ,
                            JoltCollisionHit* outHits, int maxHits, float penetrationTolerance)
{
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);
	const Shape* s = static_cast<const Shape*>(shape);

	// Get narrow phase query interface
	const NarrowPhaseQuery& query = ps->GetNarrowPhaseQuery();

	// Create filter adapters (query shape acts as MOVING layer)
	BroadPhaseLayerFilterAdapter bpFilter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter objFilter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Create collector to gather all hits
	AllHitsCollector collector(outHits, maxHits);

	// Setup collision settings
	CollideShapeSettings settings;
	settings.mPenetrationTolerance = penetrationTolerance;

	// Perform collision query
	query.CollideShape(
		s,
		Vec3::sReplicate(1.0f),  // Scale
		RMat44::sTranslation(RVec3(posX, posY, posZ)),  // Transform (position, no rotation)
		settings,
		RVec3::sZero(),  // Base offset
		collector,
		bpFilter,
		objFilter
	);

	return collector.GetNumHits();
}

// Raycast: Closest hit collector
class ClosestRayHitCollector : public CastRayCollector
{
public:
	ClosestRayHitCollector() : m_hasHit(false), m_closestFraction(1.0f) {}

	virtual void AddHit(const RayCastResult& inResult) override
	{
		if (inResult.mFraction < m_closestFraction)
		{
			m_hasHit = true;
			m_closestFraction = inResult.mFraction;
			m_closestHit = inResult;
		}
	}

	bool HasHit() const { return m_hasHit; }
	const RayCastResult& GetClosestHit() const { return m_closestHit; }

private:
	bool m_hasHit;
	float m_closestFraction;
	RayCastResult m_closestHit;
};

// Raycast: All hits collector
class AllRayHitsCollector : public CastRayCollector
{
public:
	AllRayHitsCollector(JoltRaycastHit* outHits, int maxHits)
		: m_outHits(outHits), m_maxHits(maxHits), m_numHits(0) {}

	virtual void AddHit(const RayCastResult& inResult) override
	{
		if (m_numHits < m_maxHits)
		{
			m_hits.push_back(inResult);
		}
	}

	void Finalize(const RRayCast& ray, PhysicsSystem* ps)
	{
		// Sort hits by distance (fraction)
		std::sort(m_hits.begin(), m_hits.end(),
			[](const RayCastResult& a, const RayCastResult& b) {
				return a.mFraction < b.mFraction;
			});

		// Get body lock interface for reading body data
		const BodyLockInterface& bodyLock = ps->GetBodyLockInterface();

		// Convert to output format
		int numToReturn = std::min(static_cast<int>(m_hits.size()), m_maxHits);
		for (int i = 0; i < numToReturn; i++)
		{
			const RayCastResult& result = m_hits[i];
			JoltRaycastHit& hit = m_outHits[i];

			// Store body ID
			BodyID* bodyIDCopy = new BodyID(result.mBodyID);
			hit.bodyID = static_cast<JoltBodyID>(bodyIDCopy);

			// Calculate hit point
			RVec3 hitPoint = ray.GetPointOnRay(result.mFraction);
			hit.hitPointX = static_cast<float>(hitPoint.GetX());
			hit.hitPointY = static_cast<float>(hitPoint.GetY());
			hit.hitPointZ = static_cast<float>(hitPoint.GetZ());

			// Get surface normal from the body
			Vec3 normal = Vec3::sZero();
			{
				BodyLockRead lock(bodyLock, result.mBodyID);
				if (lock.Succeeded())
				{
					const Body& body = lock.GetBody();
					normal = body.GetWorldSpaceSurfaceNormal(result.mSubShapeID2, Vec3(ray.GetPointOnRay(result.mFraction)));
				}
			}
			hit.normalX = normal.GetX();
			hit.normalY = normal.GetY();
			hit.normalZ = normal.GetZ();

			// Store fraction
			hit.fraction = result.mFraction;
		}
		m_numHits = numToReturn;
	}

	int GetNumHits() const { return m_numHits; }

private:
	JoltRaycastHit* m_outHits;
	int m_maxHits;
	int m_numHits;
	std::vector<RayCastResult> m_hits;
};

int JoltCastRay(JoltPhysicsSystem system,
                float originX, float originY, float originZ,
                float directionX, float directionY, float directionZ,
                JoltRaycastHit* outHit)
{
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);

	// Get narrow phase query interface
	const NarrowPhaseQuery& query = ps->GetNarrowPhaseQuery();

	// Create the ray
	RRayCast ray;
	ray.mOrigin = RVec3(originX, originY, originZ);
	ray.mDirection = Vec3(directionX, directionY, directionZ);

	// Create filter adapters (ray acts as MOVING layer)
	BroadPhaseLayerFilterAdapter bpFilter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter objFilter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Create collector for closest hit
	ClosestRayHitCollector collector;

	// Create raycast settings
	RayCastSettings settings;

	// Perform raycast
	query.CastRay(
		ray,
		settings,
		collector,
		bpFilter,
		objFilter
	);

	// Store result if hit and outHit is provided
	if (collector.HasHit() && outHit != nullptr)
	{
		const RayCastResult& result = collector.GetClosestHit();

		// Store body ID
		BodyID* bodyIDCopy = new BodyID(result.mBodyID);
		outHit->bodyID = static_cast<JoltBodyID>(bodyIDCopy);

		// Calculate hit point
		RVec3 hitPoint = ray.GetPointOnRay(result.mFraction);
		outHit->hitPointX = static_cast<float>(hitPoint.GetX());
		outHit->hitPointY = static_cast<float>(hitPoint.GetY());
		outHit->hitPointZ = static_cast<float>(hitPoint.GetZ());

		// Get surface normal from the body
		Vec3 normal = Vec3::sZero();
		{
			const BodyLockInterface& bodyLock = ps->GetBodyLockInterface();
			BodyLockRead lock(bodyLock, result.mBodyID);
			if (lock.Succeeded())
			{
				const Body& body = lock.GetBody();
				normal = body.GetWorldSpaceSurfaceNormal(result.mSubShapeID2, Vec3(ray.GetPointOnRay(result.mFraction)));
			}
		}
		outHit->normalX = normal.GetX();
		outHit->normalY = normal.GetY();
		outHit->normalZ = normal.GetZ();

		// Store fraction
		outHit->fraction = result.mFraction;
	}

	return collector.HasHit() ? 1 : 0;
}

int JoltCastRayGetHits(JoltPhysicsSystem system,
                       float originX, float originY, float originZ,
                       float directionX, float directionY, float directionZ,
                       JoltRaycastHit* outHits, int maxHits)
{
	PhysicsSystemWrapper* wrapper = static_cast<PhysicsSystemWrapper*>(system);
	PhysicsSystem* ps = GetPhysicsSystem(wrapper);

	// Get narrow phase query interface
	const NarrowPhaseQuery& query = ps->GetNarrowPhaseQuery();

	// Create the ray
	RRayCast ray;
	ray.mOrigin = RVec3(originX, originY, originZ);
	ray.mDirection = Vec3(directionX, directionY, directionZ);

	// Create filter adapters (ray acts as MOVING layer)
	BroadPhaseLayerFilterAdapter bpFilter(GetObjectVsBroadPhaseLayerFilter(wrapper), Layers::MOVING);
	ObjectLayerFilterAdapter objFilter(GetObjectLayerPairFilter(wrapper), Layers::MOVING);

	// Create collector for all hits
	AllRayHitsCollector collector(outHits, maxHits);

	// Create raycast settings
	RayCastSettings settings;

	// Perform raycast
	query.CastRay(
		ray,
		settings,
		collector,
		bpFilter,
		objFilter
	);

	// Finalize results (sorts and converts)
	collector.Finalize(ray, ps);

	return collector.GetNumHits();
}
