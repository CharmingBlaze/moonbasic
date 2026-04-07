/*
 * Jolt Physics C Wrapper - Collision Queries
 *
 * Handles shape overlap tests and collision detection queries.
 */

#ifndef JOLT_WRAPPER_QUERY_H
#define JOLT_WRAPPER_QUERY_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types (defined in other headers)
typedef void* JoltPhysicsSystem;
typedef void* JoltShape;
typedef void* JoltBodyID;

// Result structure for collision hits
typedef struct {
    JoltBodyID bodyID;
    float contactPointX;
    float contactPointY;
    float contactPointZ;
    float penetrationDepth;
} JoltCollisionHit;

// Result structure for raycast hits
typedef struct {
    JoltBodyID bodyID;      // The body that was hit (NULL if no hit)
    float hitPointX;        // Hit position in world space
    float hitPointY;
    float hitPointZ;
    float normalX;          // Surface normal at hit point
    float normalY;
    float normalZ;
    float fraction;         // Fraction along the ray where hit occurred [0, 1]
} JoltRaycastHit;

// Check if a shape at a position collides with anything in the physics system
// Returns 1 if collision detected, 0 if no collision
// penetrationTolerance: distance threshold for collision detection (use 0 for default)
int JoltCollideShape(JoltPhysicsSystem system, JoltShape shape,
                     float posX, float posY, float posZ, float penetrationTolerance);

// Get all collision hits for a shape at a position
// outHits: array to store results (allocated by caller)
// maxHits: maximum number of hits to return
// penetrationTolerance: distance threshold for collision detection (use 0 for default)
// Returns: actual number of hits found (may be less than maxHits)
int JoltCollideShapeGetHits(JoltPhysicsSystem system, JoltShape shape,
                            float posX, float posY, float posZ,
                            JoltCollisionHit* outHits, int maxHits, float penetrationTolerance);

// Cast a ray and check if it hits anything
// Returns 1 if hit detected, 0 if no hit
// outHit: pointer to store the closest hit result (can be NULL if you only need hit/no-hit)
int JoltCastRay(JoltPhysicsSystem system,
                float originX, float originY, float originZ,
                float directionX, float directionY, float directionZ,
                JoltRaycastHit* outHit);

// Cast a ray and get all hits along the ray (sorted by distance)
// outHits: array to store results (allocated by caller)
// maxHits: maximum number of hits to return
// Returns: actual number of hits found (may be less than maxHits)
int JoltCastRayGetHits(JoltPhysicsSystem system,
                       float originX, float originY, float originZ,
                       float directionX, float directionY, float directionZ,
                       JoltRaycastHit* outHits, int maxHits);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_QUERY_H
