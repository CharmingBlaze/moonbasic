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
		"ENTITY.LOADMESH", "ENTITY.LOADANIMATEDMESH",
		"ENTITY.SETPOSITION", "ENTITY.GETPOSITION",
		"ENTITY.POSITIONENTITY", "ENTITY.ROTATEENTITY", "ENTITY.TURNENTITY",
		"ENTITY.MOVEENTITY", "ENTITY.TRANSLATEENTITY", "ENTITY.GRAVITY",
		"ENTITY.MOVE", "ENTITY.TRANSLATE", "ENTITY.ROTATE", "ENTITY.SCALE", "ENTITY.COLOR",
		"ENTITY.ENTITYX", "ENTITY.ENTITYY", "ENTITY.ENTITYZ",
		"ENTITY.ENTITYPITCH", "ENTITY.ENTITYYAW", "ENTITY.ENTITYROLL",
		"ENTITY.PARENT", "ENTITY.PARENTCLEAR",
		"ENTITY.VISIBLE", "EntityVisible",
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
		"ENTITY.FLOOR", "ENTITY.UPDATE", "ENTITY.DRAWALL",
		"DrawEntities", "MoveEntity", "TranslateEntity", "TFormVector", "EntityHitsType", "EntityGrounded", "EntityMoveCameraRelative", "ENTITY.MOVECAMERARELATIVE", "CreatePivot", "CreateCube", "CreateSphere", "CreateCylinder", "CreateCamera",
		"EntityPBR", "EntityNormalMap", "EntityEmission",
		"EntityMass", "EntityFriction", "EntityRestitution", "ApplyEntityImpulse",
		"CameraSmoothFollow", "CreateVehicle", "AddWheel",
		"ENTITY.POINTENTITY", "ENTITY.ALIGNTOVECTOR",
		"ENTITY.ANIMATE", "ENTITY.SETANIMTIME", "ENTITY.ANIMTIME", "ENTITY.ANIMLENGTH",
		"ENTITY.EXTRACTANIMSEQ", "ENTITY.SETANIMINDEX", "ENTITY.FINDBONE",
		"LoadMesh", "LoadAnimMesh", "Animate", "SetAnimTime", "EntityAnimTime", "FindBone", "ExtractAnimSeq",
		"CreateBrush", "BrushTexture", "BrushFX", "BrushShininess", "PaintEntity", "EntityShadow",
		"LoadBrush", "FreeBrush", "BrushColor", "BrushAlpha", "BrushBlend",
		"GetEntityBrush", "PaintSurface", "GetSurfaceBrush",
		"EmitSound", "CreateSurface", "AddVertex", "AddTriangle", "UpdateMesh",
		"VertexX", "VertexY", "VertexZ",
		"ENTITY.CREATESURFACE", "ENTITY.ADDVERTEX", "ENTITY.ADDTRIANGLE", "ENTITY.UPDATEMESH",
		"ENTITY.VERTEXX", "ENTITY.VERTEXY", "ENTITY.VERTEXZ",
		"ENTITY.HIDE", "ENTITY.SHOW", "ENTITY.FREE", "ENTITY.COPY", "ENTITY.SETNAME", "ENTITY.FIND",
		"ENTITY.MOVERELATIVE", "ENTITY.APPLYGRAVITY", "ENTITY.GROUNDED",
		"ENTITY.SETMASS", "ENTITY.SETFRICTION", "ENTITY.SETBOUNCE",
		"ENTITY.GROUPCREATE", "ENTITY.GROUPADD", "ENTITY.GROUPREMOVE",
		"ENTITY.ENTITIESINGROUP", "ENTITY.ENTITIESINRADIUS", "ENTITY.ENTITIESINBOX",
		"ENTITY.CLEARSCENE", "ENTITY.SAVESCENE", "ENTITY.LOADSCENE",
		"SCENE.CLEARSCENE", "SCENE.SAVESCENE", "SCENE.LOADSCENE",
		"CAMERA.FOLLOWENTITY", "CAMERA.ORBITENTITY", "CAMERA.SETTARGETENTITY", "CAMERA.CAMERAFOLLOW",
		"ENTITY.LINKPHYSBUFFER", "ENTITY.CLEARPHYSBUFFER",
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
		"EntityVisible", "EntityDistance", "EntityInView",
		"EntityType", "EntityRadius", "EntityBox", "ResetEntity",
		"ApplyEntityForce", "ApplyEntityTorque",
		"CollisionX", "CollisionY", "CollisionZ", "CollisionNX", "CollisionNY", "CollisionNZ",
		"Animate", "SetAnimTime", "EntityAnimTime", "ExtractAnimSeq", "AnimLength",
		"CreateTerrain", "CreateMirror", "EntityAutoFade",
		"UpdateNormals", "FlipMesh", "FitMesh",
		"VertexNX", "VertexNY", "VertexNZ", "VertexU", "VertexV", "CountVertices", "CountTriangles",
		"LinePick", "CameraPick", "PickedX", "PickedY", "PickedZ", "PickedNX", "PickedNY", "PickedNZ",
		"PickedEntity", "PickedDistance", "PickedSurface", "PickedTriangle",
	}
	for _, n := range names {
		r.Register(n, "entity", stub(n))
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
