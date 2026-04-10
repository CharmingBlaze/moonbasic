// Type tags for handle safety (RULE 3 — unique TypeTag per resource class).
// Every HeapObject.TypeTag() returns one of these so [Cast] fails clearly on wrong-type handles.
// New types: append exactly one tag at the end of the iota (never reorder).
package heap

// MaxSlots is the maximum number of distinct heap slots (0 is invalid).
const MaxSlots = 65535

const (
	TagNone uint16 = iota
	TagInstance
	TagArray
	TagSprite
	TagTexture
	TagFont
	TagCamera
	TagFile
	TagJSON
	TagHost
	TagPeer
	TagEvent
	TagPhysicsBody
	TagPhysicsBuilder
	TagCharController
	TagAutomationList
	TagImage
	TagMesh
	TagMaterial
	TagModel
	TagShader
	TagMatrix
	TagVec2
	TagVec3
	TagRay
	TagBBox
	TagBSphere
	TagAudioStream
	TagWave
	TagSound
	TagMusic
	TagColor
	TagMem
	TagRng
	TagStringList
	TagPhysics2D
	TagBody2D
	TagLight
	TagInstancedModel
	TagLODModel
	TagParticle
	TagTilemap
	TagAtlas
	TagCamera2D
	TagLight2D
	TagPool
	TagTween
	TagComputeShader
	TagShaderBuffer
	TagDecal
	TagNav
	TagPath
	TagNavAgent
	TagSteerGroup
	TagBTree
	TagLobby
	TagQuaternion
	TagGameTimer
	TagGameStopwatch
	TagGameTimerSim
	TagCSV
	TagTable
	TagDB
	TagDBStmt
	TagDBRows
	TagDBTx
	TagTerrain
	TagWater
	TagSky
	TagCloud
	TagWeather
	TagScatterSet
	TagProp
	TagBiome
	TagBiomeMap
	TagNoise
	TagSpriteGroup
	TagSpriteLayer
	TagSpriteBatch
	TagSpriteUI
	TagParticle2D
	TagRenderTexture
	TagEntityRef // Blitz-style entity id wrapper for handle.Method dot syntax
	TagDrawPrim3D
	TagDrawPrim2D
	TagTextDraw
	TagTextDrawEx
	TagTextureDraw
	TagImageFacade
	TagInputFacade
	TagMoverFacade
	TagJoint2D
	TagNetPacket
	TagPlayer2D // 2D XZ mover with optional axis-aligned bounds (English-style helpers)
	TagBrush       // Blitz-style material bundle (color, texture slot, FX flags) for PaintEntity
	TagMeshBuilder // CPU vertex/index buffer for procedural ENTITY.CREATEMESH / UpdateMesh
	TagTacticalGrid // GRID.* logical cell data + optional terrain projection
	TagImageSequence // IMAGE.LOADSEQUENCE — packed list of IMAGE handles for ENTITY.SETANIMATION
)
