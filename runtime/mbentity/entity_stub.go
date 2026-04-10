//go:build !cgo && !windows

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Register implements runtime.Module — stubs when Raylib is unavailable.
func (m *Module) Register(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: ENTITY.* requires CGO (Raylib)", name)
		}
	}
	names := []string{
		"CUBE", "SPHERE",
		"ENTITY.CREATE", "ENTITY.CREATEENTITY", "ENTITY.CREATEBOX", "ENTITY.CREATECUBE",
		"ENTITY.CREATESPHERE", "ENTITY.CREATECYLINDER", "ENTITY.CREATEPLANE", "ENTITY.CREATEMESH",
		"ENTITY.LOADMESH", "ENTITY.LOAD", "ENTITY.LOADANIMATEDMESH",
		"ENTITY.SETPOSITION", "ENTITY.POSITION", "ENTITY.GETPOSITION",
		"ENTITY.POSITIONENTITY", "ENTITY.SETROTATION", "ENTITY.ROTATEENTITY", "ENTITY.TURNENTITY",
		"ENTITY.MOVEENTITY", "ENTITY.TRANSLATEENTITY", "ENTITY.GRAVITY",
		"ENTITY.MOVE", "ENTITY.TRANSLATE", "ENTITY.ROTATE", "ENTITY.TURN", "ENTITY.SCALE", "ENTITY.COLOR",
		"ENTITY.ENTITYX", "ENTITY.ENTITYY", "ENTITY.ENTITYZ",
		"ENTITY.ENTITYPITCH", "ENTITY.ENTITYYAW", "ENTITY.ENTITYROLL",
		"ENTITY.PARENT", "ENTITY.PARENTCLEAR", "ENTITY.UNPARENT",
		"ENTITY.VISIBLE", "ENTITY.SETVISIBLE", "EntityVisible",
		"ENTITY.COUNTCHILDREN", "ENTITY.GETCHILD", "ENTITY.FINDCHILD",
		"ENTITY.TFORMPOINT", "ENTITY.TFORMVECTOR",
		"ENTITY.DELTAX", "ENTITY.DELTAY", "ENTITY.DELTAZ",
		"ENTITY.MATRIXELEMENT", "ENTITY.INVIEW",
		"LOADSPRITE", "ENTITY.LOADSPRITE", "ENTITY.CREATESPRITE",
		"SCALESPRITE", "SPRITEMODE", "ENTITY.SPRITEVIEWMODE", "SPRITEVIEWMODE",
		"ENTITY.RADIUS", "ENTITY.BOX",
		"ENTITY.ALPHA", "ENTITY.SHININESS", "ENTITY.TEXTURE", "ENTITY.FX", "ENTITY.BLEND", "ENTITY.ORDER",
		"ENTITY.TYPE", "ENTITY.COLLIDE",
		"ENTITY.COLLIDED", "ENTITY.COLLISIONOTHER",
		"ENTITY.COLLISIONX", "ENTITY.COLLISIONY", "ENTITY.COLLISIONZ",
		"ENTITY.COLLISIONNX", "ENTITY.COLLISIONNY", "ENTITY.COLLISIONNZ",
		"ENTITY.DISTANCE",
		"ENTITY.SETGRAVITY", "ENTITY.JUMP", "ENTITY.VELOCITY", "ENTITY.ADDFORCE",
		"ENTITY.SLIDE", "ENTITY.PICK", "ENTITY.PICKMODE",
		"ENTITY.FLOOR", "ENTITY.UPDATE", "ENTITY.DRAWALL", "ENTITY.DRAW",
		"DrawEntities", "DrawEntity", "MoveEntity", "TranslateEntity", "TFormVector", "EntityHitsType", "EntityGrounded", "EntityMoveCameraRelative", "ENTITY.MOVECAMERARELATIVE", "CreatePivot", "CreateCube", "CreateSphere", "CreateCylinder", "CreateCamera",
		"EntityPBR", "EntityNormalMap", "EntityEmission",
		"EntityMass", "EntityFriction", "EntityRestitution", "ApplyEntityImpulse",
		"CameraSmoothFollow", "CreateVehicle", "AddWheel",
		"ENTITY.POINTENTITY", "ENTITY.LOOKAT", "ENTITY.ALIGNTOVECTOR",
		"ENTITY.ANIMATE", "ENTITY.SETANIMTIME", "ENTITY.ANIMTIME", "ENTITY.ANIMLENGTH",
		"ENTITY.EXTRACTANIMSEQ", "ENTITY.SETANIMINDEX", "ENTITY.ANIMCOUNT", "ENTITY.ANIMINDEX", "ENTITY.FINDBONE",
		"ENTITY.LOADANIMATIONS", "ENTITY.PLAY", "ENTITY.PLAYNAME", "ENTITY.STOPANIM", "ENTITY.SETANIMFRAME",
		"ENTITY.SETANIMSPEED", "ENTITY.SETANIMLOOP", "ENTITY.ISPLAYING", "ENTITY.CROSSFADE", "ENTITY.TRANSITION",
		"ENTITY.GETBONEPOS", "ENTITY.GETBONEROT", "ENTITY.SETTEXTUREMAP", "MATERIAL.BULKASSIGN", "ENTITY.GETMETADATA", "ENTITY.SETSHADER", "ENTITY.GETBOUNDS", "ENTITY.RAYHIT", "ENTITY.POINTAT",
		"ENTITY.ANIMNAME$", "ENTITY.CURRENTANIM$",
		"LEVEL.SETROOT", "LEVEL.LOAD", "LEVEL.PRELOAD", "LEVEL.FINDENTITY", "LEVEL.GETMARKER", "LEVEL.GETSPAWN", "LEVEL.SHOWLAYER",
		"LEVEL.BINDSCRIPT", "LEVEL.MATCHSCRIPTBIND", "LEVEL.LOADSKYBOX", "LEVEL.OPTIMIZE",
		"LEVEL.APPLYPHYSICS", "LEVEL.SYNCLIGHTS", "PHYSICS.AUTOCREATE",
		"ENTITY.SETSTATIC", "ENTITY.SETTRIGGER", "ENTITY.INSTANCE", "TRIGGER.CREATEFROMENTITY",
		"LoadMesh", "LoadAnimMesh", "Animate", "SetAnimTime", "EntityAnimTime", "FindBone", "ExtractAnimSeq",
		"CreateBrush", "BrushTexture", "BrushFX", "BrushShininess", "PaintEntity", "EntityShadow",
		"LoadBrush", "FreeBrush", "BrushColor", "BrushAlpha", "BrushBlend",
		"GetEntityBrush", "PaintSurface", "GetSurfaceBrush",
		"EmitSound", "CreateSurface", "AddVertex", "AddTriangle", "UpdateMesh",
		"VertexX", "VertexY", "VertexZ",
		"ENTITY.CREATESURFACE", "ENTITY.ADDVERTEX", "ENTITY.ADDTRIANGLE", "ENTITY.UPDATEMESH",
		"ENTITY.VERTEXX", "ENTITY.VERTEXY", "ENTITY.VERTEXZ",
		"ENTITY.HIDE", "ENTITY.SHOW", "ENTITY.FREE", "ENTITY.COPY", "ENTITY.INSTANCEGRID", "ENTITY.SETNAME", "ENTITY.FIND",
		"ENTITY.MOVERELATIVE", "ENTITY.APPLYGRAVITY", "ENTITY.GROUNDED",
		"ENTITY.SETMASS", "ENTITY.SETFRICTION", "ENTITY.SETBOUNCE",
		"ENTITY.GROUPCREATE", "ENTITY.GROUPADD", "ENTITY.GROUPREMOVE",
		"ENTITY.ENTITIESINGROUP", "ENTITY.ENTITIESINRADIUS", "ENTITY.ENTITIESINBOX",
		"ENTITY.CLEARSCENE", "ENTITY.SAVESCENE", "ENTITY.LOADSCENE",
		"SCENE.CLEARSCENE", "SCENE.SAVESCENE", "SCENE.LOADSCENE",
		"CAMERA.FOLLOWENTITY", "CAMERA.ORBITENTITY", "CAMERA.SETTARGETENTITY", "CAMERA.CAMERAFOLLOW",
		"ENTITY.LINKPHYSBUFFER", "ENTITY.CLEARPHYSBUFFER",
		"ENTITY.SETCOLLISIONGROUP", "EntitySetCollisionGroup",
		"ENTITY.CHECKCOLLISION", "EntityCheckCollision",
		"ENTITY.RAYCAST", "EntityRaycast",
		"ENTITY.GETGROUNDNORMAL", "EntityGetGroundNormal",
		"ENTITY.APPLYIMPULSE", "EntityApplyImpulse",
		"ENTITY.CANSEE", "EntityCanSee",
		"ENTITY.GETCLOSESTWITHTAG", "EntityGetClosestWithTag",
		"ENTITY.PUSHOUTOFGEOMETRY", "EntityPushOutOfGeometry",
		"ENTITY.INFRUSTUM", "EntityInFrustum",
		"CHECK.INVIEW",
		"ENTITY.ISSUBMERGED",
		"ENTITY.LINEOFSIGHT", "EntityLineOfSight",
		"ENTITY.GETOVERLAPCOUNT", "EntityGetOverlapCount",
		"ENTITY.ANIMATETOWARD", "EntityAnimateToward",
		"ENTITY.COLLISIONLAYER", "EntityCollisionLayer",
		"EntityCollided", "ENTITYPHYSICSTOUCH",
		"PhysicsCollisionNX", "PhysicsCollisionNY", "PhysicsCollisionNZ",
		"PhysicsCollisionPX", "PhysicsCollisionPY", "PhysicsCollisionPZ", "PhysicsCollisionY", "PhysicsCollisionForce",
		"PhysicsContactCount",
		"CountCollisions",
		"ENTITY.CREATECONE", "CreateCone",
		"GetParent", "EntityParent", "FindChild", "GetChild", "CountChildren",
		"CopyEntity", "FreeEntity", "EntityPick", "Collisions", "CollisionEntity",
		"EntityColor", "EntityAlpha", "EntityShininess", "EntityTexture",
		"EntityName", "NameEntity",
		"PositionEntity", "RotateEntity", "TurnEntity", "ScaleEntity",
		"EntityX", "EntityY", "EntityZ", "EntityPitch", "EntityYaw", "EntityRoll",
		"EntityVisible", "EntityDistance", "Entity.GetDistance", "EntityInView",
		"Entity.IsType", "Entity.SendMessage", "Entity.PollMessage", "Entity.FindByProperty",
		"ENTITY.GETDISTANCE", "ENTITY.ISTYPE", "ENTITY.HASTAG", "EntityHasTag", "ENTITY.SENDMESSAGE", "ENTITY.POLLMESSAGE", "ENTITY.FINDBYPROPERTY",
		"EntityType", "EntityRadius", "EntityBox", "ResetEntity",
		"ApplyEntityForce", "ApplyEntityTorque",
		"CollisionX", "CollisionY", "CollisionZ", "CollisionNX", "CollisionNY", "CollisionNZ",
		"Animate", "SetAnimTime", "EntityAnimTime", "ExtractAnimSeq", "AnimLength",
		"CreateTerrain", "CreateMirror", "EntityAutoFade",
		"UpdateNormals", "FlipMesh", "FitMesh",
		"VertexNX", "VertexNY", "VertexNZ", "VertexU", "VertexV", "CountVertices", "CountTriangles",
		"LinePick", "CameraPick", "PickedX", "PickedY", "PickedZ", "PickedNX", "PickedNY", "PickedNZ",
		"PickedEntity", "PickedDistance", "PickedSurface", "PickedTriangle",
		"TERRAIN.SNAPY", "TERRAIN.PLACE", "ENTITY.CLAMPTOTERRAIN", "ENTITY.GETXZ", "ENTITY.FREEENTITIES", "FREEENTITIES",
	}
	for _, n := range names {
		r.Register(n, "entity", stub(n))
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
