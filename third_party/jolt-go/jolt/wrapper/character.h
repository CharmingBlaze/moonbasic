/*
 * Jolt Physics C Wrapper - Character Virtual
 *
 * Handles virtual character controller for player/NPC movement.
 */

#ifndef JOLT_WRAPPER_CHARACTER_H
#define JOLT_WRAPPER_CHARACTER_H

#ifdef __cplusplus
extern "C" {
#endif

// Opaque pointer types
typedef void* JoltCharacterVirtual;
typedef void* JoltPhysicsSystem;
typedef void* JoltShape;
typedef void* JoltBodyID;

// Ground state enum (matches Jolt's EGroundState)
typedef enum {
    JoltGroundStateOnGround = 0,       // Character is on the ground and can move freely
    JoltGroundStateOnSteepGround = 1,  // On a slope too steep to climb
    JoltGroundStateNotSupported = 2,   // Touching object but not supported (should fall)
    JoltGroundStateInAir = 3           // In the air, not touching anything
} JoltGroundState;

// Back face mode enum (matches Jolt's EBackFaceMode)
typedef enum {
    JoltBackFaceModeIgnore = 0,        // Ignore all back facing surfaces
    JoltBackFaceModeCollide = 1        // Collide with back facing surfaces
} JoltBackFaceMode;

// Character contact settings structure
typedef struct {
    int canPushCharacter;        // bool as int (0 or 1)
    int canReceiveImpulses;      // bool as int (0 or 1)
} JoltCharacterContactSettings;

// Character contact structure (subset of Jolt's CharacterVirtual::Contact)
typedef struct {
    float positionX, positionY, positionZ;              // Position where the character makes contact
    float linearVelocityX, linearVelocityY, linearVelocityZ;  // Velocity of the contact point
    float contactNormalX, contactNormalY, contactNormalZ;     // Contact normal, pointing towards the character
    float surfaceNormalX, surfaceNormalY, surfaceNormalZ;     // Surface normal of the contact
    float distance;                                     // Distance to contact (<= 0 means actual contact, > 0 means predictive)
    float fraction;                                     // Fraction along the path where this contact takes place
    JoltBodyID bodyB;                                   // ID of body we're colliding with
    unsigned long long userData;                        // User data of B
    int isSensorB;                                      // If B is a sensor (bool as int)
    int hadCollision;                                   // If the character actually collided (bool as int)
    int wasDiscarded;                                   // If contact was discarded (bool as int)
    int canPushCharacter;                               // When true, velocity can push character (bool as int)
} JoltCharacterContact;

// Character virtual settings structure
typedef struct {
    JoltShape shape;
    float upX, upY, upZ;
    float maxSlopeAngle;
    float mass;
    float maxStrength;
    float shapeOffsetX, shapeOffsetY, shapeOffsetZ;
    JoltBackFaceMode backFaceMode;
    float predictiveContactDistance;
    unsigned int maxCollisionIterations;
    unsigned int maxConstraintIterations;
    float minTimeRemaining;
    float collisionTolerance;
    float characterPadding;
    unsigned int maxNumHits;
    float hitReductionCosMaxAngle;
    float penetrationRecoverySpeed;
    int enhancedInternalEdgeRemoval;  // bool as int (0 or 1)
} JoltCharacterVirtualSettings;

// Create a new virtual character with settings at initial position (x, y, z)
JoltCharacterVirtual JoltCreateCharacterVirtual(JoltPhysicsSystem system,
                                              const JoltCharacterVirtualSettings* settings,
                                              float x, float y, float z);

// Destroy a virtual character
void JoltDestroyCharacterVirtual(JoltCharacterVirtual character);

// Update virtual character (basic update - moves character according to velocity and handles collision)
// gravityX/Y/Z: gravity vector applied when character stands on another object
void JoltCharacterVirtualUpdate(JoltCharacterVirtual character,
                                JoltPhysicsSystem system,
                                float deltaTime,
                                float gravityX, float gravityY, float gravityZ);

// Update virtual character with extended update (combines Update, StickToFloor, WalkStairs)
// gravityX/Y/Z: gravity vector applied when character stands on another object
void JoltCharacterVirtualExtendedUpdate(JoltCharacterVirtual character,
                                        JoltPhysicsSystem system,
                                        float deltaTime,
                                        float gravityX, float gravityY, float gravityZ);

// Set the linear velocity of a virtual character
void JoltCharacterVirtualSetLinearVelocity(JoltCharacterVirtual character,
                                           float x, float y, float z);

// Get the linear velocity of a virtual character
void JoltCharacterVirtualGetLinearVelocity(const JoltCharacterVirtual character,
                                           float* x, float* y, float* z);

// Get the ground velocity of a virtual character
void JoltCharacterVirtualGetGroundVelocity(const JoltCharacterVirtual character,
                                           float* x, float* y, float* z);

// Set the position of a virtual character
void JoltCharacterVirtualSetPosition(JoltCharacterVirtual character,
                                     float x, float y, float z);

// Get the position of a virtual character
void JoltCharacterVirtualGetPosition(const JoltCharacterVirtual character,
                                     float* x, float* y, float* z);

// Get the ground state of a virtual character
JoltGroundState JoltCharacterVirtualGetGroundState(const JoltCharacterVirtual character);

// Check if character is supported (on ground or steep ground)
int JoltCharacterVirtualIsSupported(const JoltCharacterVirtual character);

// Set the shape of a virtual character
// shape: new collision shape for the character
// maxPenetrationDepth: maximum allowed penetration (typically 0.1f)
// system: physics system reference
void JoltCharacterVirtualSetShape(JoltCharacterVirtual character,
                                  JoltShape shape,
                                  float maxPenetrationDepth,
                                  JoltPhysicsSystem system);

// Get the shape of a virtual character
JoltShape JoltCharacterVirtualGetShape(const JoltCharacterVirtual character);

// Get the normal of the ground surface the character is standing on
void JoltCharacterVirtualGetGroundNormal(const JoltCharacterVirtual character,
                                         float* x, float* y, float* z);

// Get the position of the ground contact point
void JoltCharacterVirtualGetGroundPosition(const JoltCharacterVirtual character,
                                           float* x, float* y, float* z);

// Get the active contacts for the character
// contacts: pointer to array to store contacts (must be pre-allocated)
// maxContacts: maximum number of contacts to return
// Returns: actual number of contacts returned
int JoltCharacterVirtualGetActiveContacts(const JoltCharacterVirtual character,
                                          JoltCharacterContact* contacts,
                                          int maxContacts);

#ifdef __cplusplus
}
#endif

#endif // JOLT_WRAPPER_CHARACTER_H
