package vm

import (
	"testing"

	"moonbasic/vm/heap"
)

func TestHandleCallDispatchZeroArgPosUsesGetters(t *testing.T) {
	k, prep, ok := handleCallDispatch(heap.TagModel, "pos", 0)
	if !ok || !prep || k != "MODEL.GETPOS" {
		t.Fatalf("model pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagModel, "pos", 3)
	if !ok || !prep || k != "MODEL.SETPOS" {
		t.Fatalf("model pos(x,y,z): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagModel, "rot", 0)
	if !ok || !prep || k != "MODEL.GETROT" {
		t.Fatalf("model rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagModel, "scale", 0)
	if !ok || !prep || k != "MODEL.GETSCALE" {
		t.Fatalf("model scale() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagInstancedModel, "pos", 0)
	if !ok || !prep || k != "INSTANCE.GETPOS" {
		t.Fatalf("instanced model pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagInstancedModel, "rot", 0)
	if !ok || !prep || k != "INSTANCE.GETROT" {
		t.Fatalf("instanced model rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagInstancedModel, "scale", 0)
	if !ok || !prep || k != "INSTANCE.GETSCALE" {
		t.Fatalf("instanced model scale() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera, "pos", 0)
	if !ok || !prep || k != "CAMERA.GETPOS" {
		t.Fatalf("camera pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera, "rot", 0)
	if !ok || !prep || k != "CAMERA.GETROT" {
		t.Fatalf("camera rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagEntityRef, "pos", 0)
	if !ok || !prep || k != "ENTITY.GETPOS" {
		t.Fatalf("entity pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagEntityRef, "rot", 0)
	if !ok || !prep || k != "ENTITY.GETROT" {
		t.Fatalf("entity rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagEntityRef, "scale", 0)
	if !ok || !prep || k != "ENTITY.GETSCALE" {
		t.Fatalf("entity scale() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagEntityRef, "col", 0)
	if !ok || !prep || k != "ENTITY.GETCOLOR" {
		t.Fatalf("entity col() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagEntityRef, "alpha", 0)
	if !ok || !prep || k != "ENTITY.GETALPHA" {
		t.Fatalf("entity alpha() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCharController, "pos", 3)
	if !ok || !prep || k != "CHARACTERREF.SETPOS" {
		t.Fatalf("character pos(x,y,z): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCharController, "rot", 0)
	if !ok || !prep || k != "CHARACTERREF.GETROT" {
		t.Fatalf("character rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagKinematicBody, "pos", 3)
	if !ok || !prep || k != "BODYREF.SETPOS" {
		t.Fatalf("kinematic body pos(x,y,z): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagKinematicBody, "pos", 0)
	if !ok || !prep || k != "BODYREF.GETPOSITION" {
		t.Fatalf("kinematic body pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagStaticBody, "rot", 0)
	if !ok || !prep || k != "BODYREF.GETROTATION" {
		t.Fatalf("static body rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagTriggerBody, "pos", 0)
	if !ok || !prep || k != "BODYREF.GETPOSITION" {
		t.Fatalf("trigger body pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagPhysicsBody, "pos", 0)
	if !ok || !prep || k != "BODY3D.GETPOS" {
		t.Fatalf("body3d pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagPhysicsBody, "rot", 0)
	if !ok || !prep || k != "BODY3D.GETROT" {
		t.Fatalf("body3d rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagPhysicsBody, "scale", 0)
	if !ok || !prep || k != "BODY3D.GETSCALE" {
		t.Fatalf("body3d scale() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagBody2D, "pos", 0)
	if !ok || !prep || k != "BODY2D.GETPOS" {
		t.Fatalf("body2d pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagBody2D, "rot", 0)
	if !ok || !prep || k != "BODY2D.GETROT" {
		t.Fatalf("body2d rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight, "pos", 0)
	if !ok || !prep || k != "LIGHT.GETPOS" {
		t.Fatalf("light pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight, "rot", 0)
	if !ok || !prep || k != "LIGHT.GETDIR" {
		t.Fatalf("light rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight, "col", 0)
	if !ok || !prep || k != "LIGHT.GETCOLOR" {
		t.Fatalf("light col() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "pos", 0)
	if !ok || !prep || k != "SPRITE.GETPOS" {
		t.Fatalf("sprite pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "rot", 0)
	if !ok || !prep || k != "SPRITE.GETROT" {
		t.Fatalf("sprite rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "scale", 0)
	if !ok || !prep || k != "SPRITE.GETSCALE" {
		t.Fatalf("sprite scale() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "col", 0)
	if !ok || !prep || k != "SPRITE.GETCOLOR" {
		t.Fatalf("sprite col() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "alpha", 0)
	if !ok || !prep || k != "SPRITE.GETALPHA" {
		t.Fatalf("sprite alpha() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "PointHit", 2)
	if !ok || !prep || k != "SPRITE.POINTHIT" {
		t.Fatalf("sprite PointHit(x,y): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "Hit", 1)
	if !ok || !prep || k != "SPRITE.HIT" {
		t.Fatalf("sprite Hit(other): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "Collide", 1)
	if !ok || !prep || k != "SPRITE.HIT" {
		t.Fatalf("sprite Collide(other): got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagSprite, "SetFrame", 1)
	if !ok || !prep || k != "SPRITE.SETFRAME" {
		t.Fatalf("sprite SetFrame: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagParticle, "pos", 0)
	if !ok || !prep || k != "PARTICLE.GETPOS" {
		t.Fatalf("particle pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagParticle, "col", 0)
	if !ok || !prep || k != "PARTICLE.GETCOLOR" {
		t.Fatalf("particle col() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagParticle, "alpha", 0)
	if !ok || !prep || k != "PARTICLE.GETALPHA" {
		t.Fatalf("particle alpha() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagNavAgent, "pos", 0)
	if !ok || !prep || k != "NAVAGENT.GETPOS" {
		t.Fatalf("navagent pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagNavAgent, "rot", 0)
	if !ok || !prep || k != "NAVAGENT.GETROT" {
		t.Fatalf("navagent rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight2D, "pos", 0)
	if !ok || !prep || k != "LIGHT2D.GETPOS" {
		t.Fatalf("light2d pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight2D, "col", 0)
	if !ok || !prep || k != "LIGHT2D.GETCOLOR" {
		t.Fatalf("light2d col() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera2D, "pos", 0)
	if !ok || !prep || k != "CAMERA2D.GETPOS" {
		t.Fatalf("camera2d pos() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera2D, "rot", 0)
	if !ok || !prep || k != "CAMERA2D.GETROTATION" {
		t.Fatalf("camera2d rot() getter: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight2D, "getPos", 0)
	if !ok || !prep || k != "LIGHT2D.GETPOS" {
		t.Fatalf("light2d getPos() explicit: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagLight2D, "getColor", 0)
	if !ok || !prep || k != "LIGHT2D.GETCOLOR" {
		t.Fatalf("light2d getColor() explicit: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera2D, "getPos", 0)
	if !ok || !prep || k != "CAMERA2D.GETPOS" {
		t.Fatalf("camera2d getPos() explicit: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera2D, "getRotation", 0)
	if !ok || !prep || k != "CAMERA2D.GETROTATION" {
		t.Fatalf("camera2d getRotation() explicit: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera2D, "getRot", 0)
	if !ok || !prep || k != "CAMERA2D.GETROTATION" {
		t.Fatalf("camera2d getRot() alias: got %q prepend=%v ok=%v", k, prep, ok)
	}
	k, prep, ok = handleCallDispatch(heap.TagCamera2D, "getMatrix", 0)
	if !ok || !prep || k != "CAMERA2D.GETMATRIX" {
		t.Fatalf("camera2d getMatrix() explicit: got %q prepend=%v ok=%v", k, prep, ok)
	}
}

// Explicit GET* names must map with prependReceiver (handleCallBuiltin often omits these for EntityRef / NavAgent / BODYREF).
func TestHandleCallDispatchTilemapWidthHeightAndSize(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagTilemap, "getWidth", 0)
	if !o || !p || k != "TILEMAP.WIDTH" {
		t.Fatalf("tilemap getWidth: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTilemap, "getHeight", 0)
	if !o || !p || k != "TILEMAP.HEIGHT" {
		t.Fatalf("tilemap getHeight: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTilemap, "size", 0)
	if !o || !p || k != "TILEMAP.WIDTH" {
		t.Fatalf("tilemap size(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchPoolGetAndFree(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagPool, "get", 0)
	if !o || !p || k != "POOL.GET" {
		t.Fatalf("pool get: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPool, "free", 0)
	if !o || !p || k != "POOL.FREE" {
		t.Fatalf("pool free: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchAtlasGetSpriteAndFree(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagAtlas, "getSprite", 0)
	if !o || !p || k != "ATLAS.GETSPRITE" {
		t.Fatalf("atlas getSprite: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAtlas, "free", 0)
	if !o || !p || k != "ATLAS.FREE" {
		t.Fatalf("atlas free: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchPathAutomationListSizeAfterSetSizeNormalize(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagPath, "size", 0)
	if !o || !p || k != "PATH.NODECOUNT" {
		t.Fatalf("path size(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAutomationList, "size", 0)
	if !o || !p || k != "EVENT.LISTCOUNT" {
		t.Fatalf("automation list size(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchJsonCsvTableSizeAfterSetSizeNormalize(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagJSON, "size", 0)
	if !o || !p || k != "JSON.LEN" {
		t.Fatalf("json size(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagCSV, "size", 0)
	if !o || !p || k != "CSV.ROWCOUNT" {
		t.Fatalf("csv size(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTable, "size", 0)
	if !o || !p || k != "TABLE.ROWCOUNT" {
		t.Fatalf("table size(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchArrayFileSizeAfterSetSizeNormalize(t *testing.T) {
	// normalizeHandleMethod maps SIZE -> SETSIZE; 0-arg .size() must still resolve to length/getsize.
	k, p, o := handleCallDispatch(heap.TagArray, "size", 0)
	if !o || !p || k != "ARRAY.LEN" {
		t.Fatalf("array size(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagFile, "size", 0)
	if !o || !p || k != "FILE.GETSIZE" {
		t.Fatalf("file size(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagFile, "pos", 0)
	if !o || !p || k != "FILE.GETPOS" {
		t.Fatalf("file pos() after SETPOS normalize: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchTextureImageSizeAndLightEnergy(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagTexture, "size", 0)
	if !o || !p || k != "TEXTURE.GETSIZE" {
		t.Fatalf("texture size() 0-arg (SETSIZE normalize): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagImage, "getWidth", 0)
	if !o || !p || k != "IMAGE.GETWIDTH" {
		t.Fatalf("image getWidth(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagImageSequence, "size", 0)
	if !o || !p || k != "IMAGE.GETSIZE" {
		t.Fatalf("image sequence size(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagImageSequence, "FREE")
	if !o || !p || k != "IMAGE.FREE" {
		t.Fatalf("image sequence FREE builtin: got %q prepend=%v ok=%v", k, p, o)
	}
	ke, pe, oe := handleCallBuiltin(heap.TagLight, "getEnergy")
	if !oe || !pe || ke != "LIGHT.GETENERGY" {
		t.Fatalf("light getEnergy builtin: got %q prepend=%v ok=%v", ke, pe, oe)
	}
}

func TestHandleCallDispatchPlayer2DPosZeroArg(t *testing.T) {
	for _, method := range []string{"pos", "position"} {
		k, p, o := handleCallDispatch(heap.TagPlayer2D, method, 0)
		if !o || !p || k != "PLAYER2D.GETPOS" {
			t.Fatalf("player2d %s() 0-arg: got %q prepend=%v ok=%v", method, k, p, o)
		}
	}
	k, p, o := handleCallDispatch(heap.TagPlayer2D, "getPos", 0)
	if !o || !p || k != "PLAYER2D.GETPOS" {
		t.Fatalf("player2d getPos(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPlayer2D, "SETPOS")
	if !o || !p || k != "PLAYER2D.SETPOS" {
		t.Fatalf("player2d SETPOS (write): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchParticleSizeZeroArg(t *testing.T) {
	for _, method := range []string{"size", "scale"} {
		k, p, o := handleCallDispatch(heap.TagParticle, method, 0)
		if !o || !p || k != "PARTICLE.GETSIZE" {
			t.Fatalf("particle %s() 0-arg: got %q prepend=%v ok=%v", method, k, p, o)
		}
	}
	k, p, o := handleCallDispatch(heap.TagParticle, "getSize", 0)
	if !o || !p || k != "PARTICLE.GETSIZE" {
		t.Fatalf("particle getSize(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallBuiltinTextureVsImageMutateOps(t *testing.T) {
	k, p, o := handleCallBuiltin(heap.TagTexture, "CROP")
	if o {
		t.Fatalf("texture CROP: expected no builtin mapping, got %q prepend=%v", k, p)
	}
	k, p, o = handleCallBuiltin(heap.TagImage, "CROP")
	if !o || !p || k != "IMAGE.CROP" {
		t.Fatalf("image CROP: got %q prepend=%v ok=%v", k, p, o)
	}
	_, _, o = handleCallBuiltin(heap.TagTexture, "DRAW")
	if o {
		t.Fatal("texture DRAW: expected no mapping (no TEXTURE.DRAW)")
	}
	_, _, o = handleCallBuiltin(heap.TagImage, "DRAW")
	if o {
		t.Fatal("image DRAW: expected no mapping (no IMAGE.DRAW; use IMAGE.DRAW* or draw commands)")
	}
	_, _, o = handleCallBuiltin(heap.TagFont, "TEXTWIDTH")
	if o {
		t.Fatal("font TEXTWIDTH: expected no mapping (use DRAW.TEXTFONTWIDTH / DRAW.TEXTWIDTH)")
	}
}

func TestHandleCallDispatchCamera3DExplicitGettersUseBuiltinFallback(t *testing.T) {
	// 0-arg: dispatch switch uses SETROT→GETROT for "rot" but not "getRot"; must fall through to handleCallBuiltin allowlist.
	tests := []struct {
		method string
		want   string
	}{
		{"getRot", "CAMERA.GETROT"},
		{"getFov", "CAMERA.GETFOV"},
		{"getUp", "CAMERA.GETUP"},
		{"getProjection", "CAMERA.GETPROJECTION"},
	}
	for _, tc := range tests {
		k, p, o := handleCallDispatch(heap.TagCamera, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("TagCamera %q: got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
}

func TestHandleCallDispatchExplicitGetterEntityNavBodyRef(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagEntityRef, "getPos", 0)
	if !o || !p || k != "ENTITY.GETPOS" {
		t.Fatalf("entity getPos: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagNavAgent, "getSpeed", 0)
	if !o || !p || k != "NAVAGENT.GETSPEED" {
		t.Fatalf("navagent getSpeed: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagNavAgent, "isAtDestination", 0)
	if !o || !p || k != "NAVAGENT.ISATDESTINATION" {
		t.Fatalf("navagent isAtDestination: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagKinematicBody, "getPosition", 0)
	if !o || !p || k != "BODYREF.GETPOSITION" {
		t.Fatalf("kinematic getPosition: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagStaticBody, "getRot", 0)
	if !o || !p || k != "BODYREF.GETROTATION" {
		t.Fatalf("static getRot: got %q prepend=%v ok=%v", k, p, o)
	}
}

// handleCallBuiltin omits AUDIOSTREAM.ISPLAYING / ISREADY and TWEEN.ISPLAYING / ISFINISHED / PROGRESS;
// handleCallDispatch must map them so 0-arg calls prepend the receiver.
func TestHandleCallBuiltinParticleSetScaleMapsToSetSize(t *testing.T) {
	// normalize maps scale → SETSCALE; must not resolve to a non-existent PARTICLE.SETSCALE.
	k, p, ok := handleCallBuiltin(heap.TagParticle, "scale")
	if !ok || !p || k != "PARTICLE.SETSIZE" {
		t.Fatalf("particle scale(...): got %q prepend=%v ok=%v want PARTICLE.SETSIZE", k, p, ok)
	}
}

func TestHandleCallDispatchParticleVelocityZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagParticle, "setVelocity", 0)
	if !o || !p || k != "PARTICLE.GETVELOCITY" {
		t.Fatalf("particle setVelocity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagParticle, "vel", 0)
	if !o || !p || k != "PARTICLE.GETVELOCITY" {
		t.Fatalf("particle vel(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchVecLenAlias(t *testing.T) {
	for _, tc := range []struct {
		tag  uint16
		want string
	}{
		{heap.TagVec2, "VEC2.LENGTH"},
		{heap.TagVec3, "VEC3.LENGTH"},
	} {
		k, p, o := handleCallDispatch(tc.tag, "len", 0)
		if !o || !p || k != tc.want {
			t.Fatalf("tag %d len(): got %q prepend=%v ok=%v want %q", tc.tag, k, p, o, tc.want)
		}
		k, p, o = handleCallBuiltin(tc.tag, "len")
		if !o || !p || k != tc.want {
			t.Fatalf("tag %d len(...): got %q prepend=%v ok=%v want %q", tc.tag, k, p, o, tc.want)
		}
	}
}

func TestHandleCallDispatchTweenLoopYoyoZeroArgAreGetters(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagTween, "setLoop", 0)
	if !o || !p || k != "TWEEN.GETLOOP" {
		t.Fatalf("tween setLoop(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTween, "loop", 0)
	if !o || !p || k != "TWEEN.GETLOOP" {
		t.Fatalf("tween loop(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTween, "yoyo", 0)
	if !o || !p || k != "TWEEN.GETYOYO" {
		t.Fatalf("tween yoyo(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTween, "getYoyo", 0)
	if !o || !p || k != "TWEEN.GETYOYO" {
		t.Fatalf("tween getYoyo(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchAudioStreamAndTweenQueryMethods(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagAudioStream, "isPlaying", 0)
	if !o || !p || k != "AUDIOSTREAM.ISPLAYING" {
		t.Fatalf("audiostream isPlaying: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "isReady", 0)
	if !o || !p || k != "AUDIOSTREAM.ISREADY" {
		t.Fatalf("audiostream isReady: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "setVolume", 0)
	if !o || !p || k != "AUDIOSTREAM.GETVOLUME" {
		t.Fatalf("audiostream setVolume(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "setPitch", 0)
	if !o || !p || k != "AUDIOSTREAM.GETPITCH" {
		t.Fatalf("audiostream setPitch(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "setPan", 0)
	if !o || !p || k != "AUDIOSTREAM.GETPAN" {
		t.Fatalf("audiostream setPan(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "getVolume", 0)
	if !o || !p || k != "AUDIOSTREAM.GETVOLUME" {
		t.Fatalf("audiostream getVolume(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "getPitch", 0)
	if !o || !p || k != "AUDIOSTREAM.GETPITCH" {
		t.Fatalf("audiostream getPitch(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagAudioStream, "getPan", 0)
	if !o || !p || k != "AUDIOSTREAM.GETPAN" {
		t.Fatalf("audiostream getPan(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTween, "isPlaying", 0)
	if !o || !p || k != "TWEEN.ISPLAYING" {
		t.Fatalf("tween isPlaying: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTween, "isFinished", 0)
	if !o || !p || k != "TWEEN.ISFINISHED" {
		t.Fatalf("tween isFinished: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagTween, "progress", 0)
	if !o || !p || k != "TWEEN.PROGRESS" {
		t.Fatalf("tween progress: got %q prepend=%v ok=%v", k, p, o)
	}
}

// normalizeHandleMethod maps "pos" → SETPOS; without these cases, 0-arg calls would fall through to RAY.SETPOS / BBOX.SET* / BSPHERE.SET*.
func TestHandleCallDispatchSoundMusicSetVolumeSetPitchNormalize(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagSound, "setVolume", 0)
	if !o || !p || k != "AUDIO.GETSOUNDVOLUME" {
		t.Fatalf("sound setVolume(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagMusic, "setVolume", 0)
	if !o || !p || k != "AUDIO.GETMUSICVOLUME" {
		t.Fatalf("music setVolume(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagSound, "setPitch", 0)
	if !o || !p || k != "AUDIO.GETSOUNDPITCH" {
		t.Fatalf("sound setPitch(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagMusic, "setPitch", 0)
	if !o || !p || k != "AUDIO.GETMUSICPITCH" {
		t.Fatalf("music setPitch(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagSound, "setPan", 0)
	if !o || !p || k != "AUDIO.GETSOUNDPAN" {
		t.Fatalf("sound setPan(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchWaterWaveGetterNotSetter(t *testing.T) {
	// handleCallBuiltin maps WAVE → WATER.SETWAVE; bare "wave" must dispatch to GETWAVESPEED for 0-arg.
	k, p, o := handleCallDispatch(heap.TagWater, "wave", 0)
	if !o || !p || k != "WATER.GETWAVESPEED" {
		t.Fatalf("water wave(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchDecalSetSizeLifetimeZeroArgAreGetters(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagDecal, "setSize", 0)
	if !o || !p || k != "DECAL.GETSIZE" {
		t.Fatalf("decal setSize(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagDecal, "setLifetime", 0)
	if !o || !p || k != "DECAL.GETLIFETIME" {
		t.Fatalf("decal setLifetime(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagDecal, "getRot", 0)
	if !o || !p || k != "DECAL.GETROT" {
		t.Fatalf("decal getRot(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchGameTimerSimSetLoopZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagGameTimerSim, "setLoop", 0)
	if !o || !p || k != "TIMER.GETLOOP" {
		t.Fatalf("GameTimerSim setLoop(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagGameTimerSim, "loop", 0)
	if !o || !p || k != "TIMER.GETLOOP" {
		t.Fatalf("GameTimerSim loop(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagGameTimerSim, "getLoop", 0)
	if !o || !p || k != "TIMER.GETLOOP" {
		t.Fatalf("GameTimerSim getLoop(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchCloudSetCoverageZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagCloud, "setCoverage", 0)
	if !o || !p || k != "CLOUD.GETCOVERAGE" {
		t.Fatalf("cloud setCoverage(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagCloud, "getCoverage", 0)
	if !o || !p || k != "CLOUD.GETCOVERAGE" {
		t.Fatalf("cloud getCoverage(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBiomeSetTempHumidityZeroArgAreGetters(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagBiome, "setTemp", 0)
	if !o || !p || k != "BIOME.GETTEMP" {
		t.Fatalf("biome setTemp(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBiome, "setHumidity", 0)
	if !o || !p || k != "BIOME.GETHUMIDITY" {
		t.Fatalf("biome setHumidity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBiome, "getTemp", 0)
	if !o || !p || k != "BIOME.GETTEMP" {
		t.Fatalf("biome getTemp(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchSkyWeatherSetNamesZeroArgAreGetters(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagSky, "setTime", 0)
	if !o || !p || k != "SKY.GETTIMEHOURS" {
		t.Fatalf("sky setTime(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagSky, "getTimeHours", 0)
	if !o || !p || k != "SKY.GETTIMEHOURS" {
		t.Fatalf("sky getTimeHours(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagWeather, "setType", 0)
	if !o || !p || k != "WEATHER.GETTYPE" {
		t.Fatalf("weather setType(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagWeather, "setCoverage", 0)
	if !o || !p || k != "WEATHER.GETCOVERAGE" {
		t.Fatalf("weather setCoverage(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchRayBBoxBSpherePosRadiusAfterNormalize(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagRay, "pos", 0)
	if !o || !p || k != "RAY.GETPOS" {
		t.Fatalf("ray pos(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagRay, "dir", 0)
	if !o || !p || k != "RAY.GETDIR" {
		t.Fatalf("ray dir(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBBox, "min", 0)
	if !o || !p || k != "BBOX.GETMIN" {
		t.Fatalf("bbox min(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBBox, "max", 0)
	if !o || !p || k != "BBOX.GETMAX" {
		t.Fatalf("bbox max(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBSphere, "pos", 0)
	if !o || !p || k != "BSPHERE.GETPOS" {
		t.Fatalf("bsphere pos(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBSphere, "radius", 0)
	if !o || !p || k != "BSPHERE.GETRADIUS" {
		t.Fatalf("bsphere radius(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallBuiltinBody2DRotSetter(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagBody2D, "rot")
	if !ok || !prep || k != "BODY2D.SETROT" {
		t.Fatalf("body2d rot(angle): got %q prepend=%v ok=%v", k, prep, ok)
	}
}

func TestHandleCallDispatchBody2DSetVelAliasesZeroArgAreGetters(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagBody2D, "setVel", 0)
	if !o || !p || k != "BODY2D.GETLINEARVELOCITY" {
		t.Fatalf("body2d setVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBody2D, "setAngularVel", 0)
	if !o || !p || k != "BODY2D.GETANGULARVELOCITY" {
		t.Fatalf("body2d setAngularVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBody2D, "setMass", 0)
	if !o || !p || k != "BODY2D.GETMASS" {
		t.Fatalf("body2d setMass(): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBody3DVelAliasZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagPhysicsBody, "vel", 0)
	if !o || !p || k != "BODY3D.GETVELOCITY" {
		t.Fatalf("body3d vel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPhysicsBody, "velocity", 0)
	if !o || !p || k != "BODY3D.GETVELOCITY" {
		t.Fatalf("body3d velocity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "vel")
	if !o || !p || k != "BODY3D.SETVELOCITY" {
		t.Fatalf("body3d vel(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "velocity")
	if !o || !p || k != "BODY3D.SETVELOCITY" {
		t.Fatalf("body3d velocity(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBody3DMassAliasZeroArgAndWritePath(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagPhysicsBody, "mass", 0)
	if !o || !p || k != "BODY3D.GETMASS" {
		t.Fatalf("body3d mass(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPhysicsBody, "setMass", 0)
	if !o || !p || k != "BODY3D.GETMASS" {
		t.Fatalf("body3d setMass(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "mass")
	if !o || !p || k != "BODY3D.SETMASS" {
		t.Fatalf("body3d mass(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "setMass")
	if !o || !p || k != "BODY3D.SETMASS" {
		t.Fatalf("body3d setMass(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBody2DAngularVelAliasZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagBody2D, "angularVel", 0)
	if !o || !p || k != "BODY2D.GETANGULARVELOCITY" {
		t.Fatalf("body2d angularVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBody2D, "angularVelocity", 0)
	if !o || !p || k != "BODY2D.GETANGULARVELOCITY" {
		t.Fatalf("body2d angularVelocity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBody2D, "angVel", 0)
	if !o || !p || k != "BODY2D.GETANGULARVELOCITY" {
		t.Fatalf("body2d angVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagBody2D, "angularVel")
	if !o || !p || k != "BODY2D.SETANGULARVELOCITY" {
		t.Fatalf("body2d angularVel(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagBody2D, "angularVelocity")
	if !o || !p || k != "BODY2D.SETANGULARVELOCITY" {
		t.Fatalf("body2d angularVelocity(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagBody2D, "angVel")
	if !o || !p || k != "BODY2D.SETANGULARVELOCITY" {
		t.Fatalf("body2d angVel(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBody2DVelocityAliasWritePath(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagBody2D, "velocity", 0)
	if !o || !p || k != "BODY2D.GETLINEARVELOCITY" {
		t.Fatalf("body2d velocity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagBody2D, "linearVel", 0)
	if !o || !p || k != "BODY2D.GETLINEARVELOCITY" {
		t.Fatalf("body2d linearVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagBody2D, "velocity")
	if !o || !p || k != "BODY2D.SETLINEARVELOCITY" {
		t.Fatalf("body2d velocity(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagBody2D, "linearVel")
	if !o || !p || k != "BODY2D.SETLINEARVELOCITY" {
		t.Fatalf("body2d linearVel(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchLightSetEnergyZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagLight, "setEnergy", 0)
	if !o || !p || k != "LIGHT.GETINTENSITY" {
		t.Fatalf("light setEnergy(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagLight, "setEnergy")
	if !o || !p || k != "LIGHT.SETINTENSITY" {
		t.Fatalf("light setEnergy(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBody3DAngularVelZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagPhysicsBody, "setAngularVel", 0)
	if !o || !p || k != "BODY3D.GETANGULARVEL" {
		t.Fatalf("body3d setAngularVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPhysicsBody, "setAngularVelocity", 0)
	if !o || !p || k != "BODY3D.GETANGULARVEL" {
		t.Fatalf("body3d setAngularVelocity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPhysicsBody, "angularVel", 0)
	if !o || !p || k != "BODY3D.GETANGULARVEL" {
		t.Fatalf("body3d angularVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPhysicsBody, "angularVelocity", 0)
	if !o || !p || k != "BODY3D.GETANGULARVEL" {
		t.Fatalf("body3d angularVelocity(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagPhysicsBody, "angVel", 0)
	if !o || !p || k != "BODY3D.GETANGULARVEL" {
		t.Fatalf("body3d angVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "setAngularVelocity")
	if !o || !p || k != "BODY3D.SETANGULARVEL" {
		t.Fatalf("body3d setAngularVelocity(...): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "angularVel")
	if !o || !p || k != "BODY3D.SETANGULARVEL" {
		t.Fatalf("body3d angularVel(...): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "angularVelocity")
	if !o || !p || k != "BODY3D.SETANGULARVEL" {
		t.Fatalf("body3d angularVelocity(...): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagPhysicsBody, "angVel")
	if !o || !p || k != "BODY3D.SETANGULARVEL" {
		t.Fatalf("body3d angVel(...): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchBody3DPhysicalPropsAliasWritePath(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"friction", "BODY3D.SETFRICTION"},
		{"bounce", "BODY3D.SETRESTITUTION"},
		{"restitution", "BODY3D.SETRESTITUTION"},
		{"damping", "BODY3D.SETDAMPING"},
		{"gravityFactor", "BODY3D.SETGRAVITYFACTOR"},
		{"ccd", "BODY3D.SETCCD"},
	}
	for _, tc := range cases {
		k, p, o := handleCallBuiltin(heap.TagPhysicsBody, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("body3d %s(...) write path: got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
}

func TestHandleCallDispatchNavAgentSetSpeedSetMaxForceZeroArgAreGetters(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagNavAgent, "setSpeed", 0)
	if !o || !p || k != "NAVAGENT.GETSPEED" {
		t.Fatalf("navagent setSpeed(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagNavAgent, "speed", 0)
	if !o || !p || k != "NAVAGENT.GETSPEED" {
		t.Fatalf("navagent speed(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagNavAgent, "setMaxForce", 0)
	if !o || !p || k != "NAVAGENT.GETMAXFORCE" {
		t.Fatalf("navagent setMaxForce(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagNavAgent, "maxForce", 0)
	if !o || !p || k != "NAVAGENT.GETMAXFORCE" {
		t.Fatalf("navagent maxForce(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "setSpeed")
	if !o || !p || k != "NAVAGENT.SETSPEED" {
		t.Fatalf("navagent setSpeed(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "setMaxForce")
	if !o || !p || k != "NAVAGENT.SETMAXFORCE" {
		t.Fatalf("navagent setMaxForce(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "maxForce")
	if !o || !p || k != "NAVAGENT.SETMAXFORCE" {
		t.Fatalf("navagent maxForce(...) write path: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchCharControllerSetAliasesZeroArgAreGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"setVelocity", "CHARACTERREF.GETVELOCITY"},
		{"vel", "CHARACTERREF.GETVELOCITY"},
		{"velocity", "CHARACTERREF.GETVELOCITY"},
		{"setAirControl", "CHARACTERREF.GETAIRCONTROL"},
		{"setGroundControl", "CHARACTERREF.GETGROUNDCONTROL"},
		{"setJumpBuffer", "CHARACTERREF.GETJUMPBUFFER"},
		{"setBounciness", "CHARACTERREF.GETBOUNCINESS"},
		{"setBounce", "CHARACTERREF.GETBOUNCINESS"},
		{"setSpeed", "CHARACTERREF.GETSPEED"},
		{"speed", "CHARACTERREF.GETSPEED"},
		{"getSpeed", "CHARACTERREF.GETSPEED"},
		{"slopeAngle", "CHARACTERREF.GETSLOPEANGLE"},
		{"getSlopeAngle", "CHARACTERREF.GETSLOPEANGLE"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagCharController, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("charcontroller %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	k, p, o := handleCallBuiltin(heap.TagCharController, "vel")
	if !o || !p || k != "CHARACTERREF.SETVELOCITY" {
		t.Fatalf("charcontroller vel(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "velocity")
	if !o || !p || k != "CHARACTERREF.SETVELOCITY" {
		t.Fatalf("charcontroller velocity(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "speed")
	if o || p || k != "" {
		t.Fatalf("charcontroller speed(...) must not map (no CHARACTERREF.SETSPEED): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "setSpeed")
	if o || p || k != "" {
		t.Fatalf("charcontroller setSpeed(...) must not map (no CHARACTERREF.SETSPEED): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "getSpeed")
	if !o || !p || k != "CHARACTERREF.GETSPEED" {
		t.Fatalf("charcontroller getSpeed(...) explicit getter: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "slopeAngle")
	if o || p || k != "" {
		t.Fatalf("charcontroller slopeAngle(...) must not map (no CHARACTERREF.SETSLOPEANGLE): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "getSlopeAngle")
	if !o || !p || k != "CHARACTERREF.GETSLOPEANGLE" {
		t.Fatalf("charcontroller getSlopeAngle(...) explicit getter: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchKinematicBodySetVelZeroArgIsGetter(t *testing.T) {
	k, p, o := handleCallDispatch(heap.TagKinematicBody, "setVel", 0)
	if !o || !p || k != "KINEMATICREF.GETVELOCITY" {
		t.Fatalf("kinematic setVel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallDispatch(heap.TagKinematicBody, "vel", 0)
	if !o || !p || k != "KINEMATICREF.GETVELOCITY" {
		t.Fatalf("kinematic vel(): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagKinematicBody, "getVelocity")
	if !o || !p || k != "KINEMATICREF.GETVELOCITY" {
		t.Fatalf("kinematic getVelocity(...): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagKinematicBody, "vel")
	if !o || !p || k != "KINEMATICREF.SETVELOCITY" {
		t.Fatalf("kinematic vel(...): got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagKinematicBody, "velocity")
	if !o || !p || k != "KINEMATICREF.SETVELOCITY" {
		t.Fatalf("kinematic velocity(...): got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchCharControllerSlopeStepdistAliasesZeroArgAreGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"setSlope", "CHARACTERREF.GETMAXSLOPE"},
		{"setMaxSlope", "CHARACTERREF.GETMAXSLOPE"},
		{"slope", "CHARACTERREF.GETMAXSLOPE"},
		{"maxSlope", "CHARACTERREF.GETMAXSLOPE"},
		{"setStep", "CHARACTERREF.GETSTEPHEIGHT"},
		{"setStepHeight", "CHARACTERREF.GETSTEPHEIGHT"},
		{"step", "CHARACTERREF.GETSTEPHEIGHT"},
		{"stepHeight", "CHARACTERREF.GETSTEPHEIGHT"},
		{"setStickDown", "CHARACTERREF.GETSNAPDISTANCE"},
		{"setSnapDistance", "CHARACTERREF.GETSNAPDISTANCE"},
		{"snap", "CHARACTERREF.GETSNAPDISTANCE"},
		{"snapDistance", "CHARACTERREF.GETSNAPDISTANCE"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagCharController, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("charcontroller %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path must still reach setters
	k, p, o := handleCallBuiltin(heap.TagCharController, "setSlope")
	if !o || !p || k != "CHARACTERREF.SETMAXSLOPE" {
		t.Fatalf("charcontroller setSlope(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "setStep")
	if !o || !p || k != "CHARACTERREF.SETSTEPHEIGHT" {
		t.Fatalf("charcontroller setStep(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCharController, "setStickDown")
	if !o || !p || k != "CHARACTERREF.SETSNAPDISTANCE" {
		t.Fatalf("charcontroller setStickDown(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchPlayer2DSetPosZeroArgIsGetter(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"setPos", "PLAYER2D.GETPOS"},
		{"pos", "PLAYER2D.GETPOS"},
		{"position", "PLAYER2D.GETPOS"},
		{"getPos", "PLAYER2D.GETPOS"},
		{"x", "PLAYER2D.GETX"},
		{"getX", "PLAYER2D.GETX"},
		{"z", "PLAYER2D.GETZ"},
		{"getZ", "PLAYER2D.GETZ"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagPlayer2D, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("player2d %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path must still reach setter
	k, p, o := handleCallBuiltin(heap.TagPlayer2D, "setPos")
	if !o || !p || k != "PLAYER2D.SETPOS" {
		t.Fatalf("player2d setPos(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchCamera2DZoomOffsetZeroArgAreGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"setZoom", "CAMERA2D.GETZOOM"},
		{"zoom", "CAMERA2D.GETZOOM"},
		{"getZoom", "CAMERA2D.GETZOOM"},
		{"setOffset", "CAMERA2D.GETOFFSET"},
		{"offset", "CAMERA2D.GETOFFSET"},
		{"getOffset", "CAMERA2D.GETOFFSET"},
		{"setTarget", "CAMERA2D.GETPOS"},
		{"target", "CAMERA2D.GETPOS"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagCamera2D, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("camera2d %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path must still reach setters
	k, p, o := handleCallBuiltin(heap.TagCamera2D, "setZoom")
	if !o || !p || k != "CAMERA2D.SETZOOM" {
		t.Fatalf("camera2d setZoom(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCamera2D, "setOffset")
	if !o || !p || k != "CAMERA2D.SETOFFSET" {
		t.Fatalf("camera2d setOffset(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagCamera2D, "setTarget")
	if !o || !p || k != "CAMERA2D.SETTARGET" {
		t.Fatalf("camera2d setTarget(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchWaterShallowDeepColorZeroArgAreGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"setShallowColor", "WATER.GETSHALLOWCOLOR"},
		{"shallowColor", "WATER.GETSHALLOWCOLOR"},
		{"setDeepColor", "WATER.GETDEEPCOLOR"},
		{"deepColor", "WATER.GETDEEPCOLOR"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagWater, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("water %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path must still reach setters
	k, p, o := handleCallBuiltin(heap.TagWater, "setShallowColor")
	if !o || !p || k != "WATER.SETSHALLOWCOLOR" {
		t.Fatalf("water setShallowColor(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagWater, "setDeepColor")
	if !o || !p || k != "WATER.SETDEEPCOLOR" {
		t.Fatalf("water setDeepColor(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchNavAgentZeroArgAreGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"setPos", "NAVAGENT.GETPOS"},
		{"pos", "NAVAGENT.GETPOS"},
		{"getPos", "NAVAGENT.GETPOS"},
		{"setRot", "NAVAGENT.GETROT"},
		{"rot", "NAVAGENT.GETROT"},
		{"getRot", "NAVAGENT.GETROT"},
		{"setSpeed", "NAVAGENT.GETSPEED"},
		{"speed", "NAVAGENT.GETSPEED"},
		{"getSpeed", "NAVAGENT.GETSPEED"},
		{"setMaxForce", "NAVAGENT.GETMAXFORCE"},
		{"getMaxForce", "NAVAGENT.GETMAXFORCE"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagNavAgent, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("navagent %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path must still reach setters
	k, p, o := handleCallBuiltin(heap.TagNavAgent, "setSpeed")
	if !o || !p || k != "NAVAGENT.SETSPEED" {
		t.Fatalf("navagent setSpeed(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "setMaxForce")
	if !o || !p || k != "NAVAGENT.SETMAXFORCE" {
		t.Fatalf("navagent setMaxForce(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "maxForce")
	if !o || !p || k != "NAVAGENT.SETMAXFORCE" {
		t.Fatalf("navagent maxForce(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	// getters via builtin (argCount>0 path)
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "getSpeed")
	if !o || !p || k != "NAVAGENT.GETSPEED" {
		t.Fatalf("navagent getSpeed(h) builtin: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagNavAgent, "getMaxForce")
	if !o || !p || k != "NAVAGENT.GETMAXFORCE" {
		t.Fatalf("navagent getMaxForce(h) builtin: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchAudioStreamVolPitchPanZeroArgAreGetters(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"volume", "AUDIOSTREAM.GETVOLUME"},
		{"setVolume", "AUDIOSTREAM.GETVOLUME"},
		{"getVolume", "AUDIOSTREAM.GETVOLUME"},
		{"pitch", "AUDIOSTREAM.GETPITCH"},
		{"setPitch", "AUDIOSTREAM.GETPITCH"},
		{"getPitch", "AUDIOSTREAM.GETPITCH"},
		{"pan", "AUDIOSTREAM.GETPAN"},
		{"setPan", "AUDIOSTREAM.GETPAN"},
		{"getPan", "AUDIOSTREAM.GETPAN"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagAudioStream, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("audiostream %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path must still reach setters
	k, p, o := handleCallBuiltin(heap.TagAudioStream, "volume")
	if !o || !p || k != "AUDIOSTREAM.SETVOLUME" {
		t.Fatalf("audiostream volume(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagAudioStream, "pitch")
	if !o || !p || k != "AUDIOSTREAM.SETPITCH" {
		t.Fatalf("audiostream pitch(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagAudioStream, "pan")
	if !o || !p || k != "AUDIOSTREAM.SETPAN" {
		t.Fatalf("audiostream pan(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchSoundMusicVolPitchZeroArgAreGetters(t *testing.T) {
	// Sound
	sCases := []struct {
		method string
		want   string
	}{
		{"volume", "AUDIO.GETSOUNDVOLUME"},
		{"setVolume", "AUDIO.GETSOUNDVOLUME"},
		{"getVolume", "AUDIO.GETSOUNDVOLUME"},
		{"pitch", "AUDIO.GETSOUNDPITCH"},
		{"setPitch", "AUDIO.GETSOUNDPITCH"},
		{"pan", "AUDIO.GETSOUNDPAN"},
		{"setPan", "AUDIO.GETSOUNDPAN"},
	}
	for _, tc := range sCases {
		k, p, o := handleCallDispatch(heap.TagSound, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("sound %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// Music
	mCases := []struct {
		method string
		want   string
	}{
		{"volume", "AUDIO.GETMUSICVOLUME"},
		{"setVolume", "AUDIO.GETMUSICVOLUME"},
		{"pitch", "AUDIO.GETMUSICPITCH"},
		{"setPitch", "AUDIO.GETMUSICPITCH"},
		{"length", "AUDIO.GETMUSICLENGTH"},
		{"time", "AUDIO.GETMUSICTIME"},
	}
	for _, tc := range mCases {
		k, p, o := handleCallDispatch(heap.TagMusic, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("music %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path: volume/pitch still reach setters
	k, p, o := handleCallBuiltin(heap.TagSound, "volume")
	if !o || !p || k != "AUDIO.SETSOUNDVOLUME" {
		t.Fatalf("sound volume(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	k, p, o = handleCallBuiltin(heap.TagMusic, "volume")
	if !o || !p || k != "AUDIO.SETMUSICVOLUME" {
		t.Fatalf("music volume(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallDispatchTimerSimLoopZeroArgIsGetter(t *testing.T) {
	cases := []struct {
		method string
		want   string
	}{
		{"loop", "TIMER.GETLOOP"},
		{"setLoop", "TIMER.GETLOOP"},
		{"getLoop", "TIMER.GETLOOP"},
	}
	for _, tc := range cases {
		k, p, o := handleCallDispatch(heap.TagGameTimerSim, tc.method, 0)
		if !o || !p || k != tc.want {
			t.Fatalf("timersim %s(): got %q prepend=%v ok=%v want %q", tc.method, k, p, o, tc.want)
		}
	}
	// write path: setLoop must reach setter
	k, p, o := handleCallBuiltin(heap.TagGameTimerSim, "setLoop")
	if !o || !p || k != "TIMER.SETLOOP" {
		t.Fatalf("timersim setLoop(...) write: got %q prepend=%v ok=%v", k, p, o)
	}
	// getLoop explicit getter via builtin
	k, p, o = handleCallBuiltin(heap.TagGameTimerSim, "getLoop")
	if !o || !p || k != "TIMER.GETLOOP" {
		t.Fatalf("timersim getLoop(...) builtin: got %q prepend=%v ok=%v", k, p, o)
	}
}

func TestHandleCallBuiltinEntitySetCollisionMeshMapsToSetStatic(t *testing.T) {
	// ENTITY.STATICENTITY never existed; correct key is ENTITY.SETSTATIC.
	for _, method := range []string{"setCollisionMesh", "collisionMesh", "setStatic", "static"} {
		k, p, o := handleCallBuiltin(heap.TagEntityRef, method)
		if !o || !p || k != "ENTITY.SETSTATIC" {
			t.Fatalf("entity %s: got %q prepend=%v ok=%v, want ENTITY.SETSTATIC", method, k, p, o)
		}
	}
}

func TestHandleCallBuiltinBody3DForceImpulseMapsToApplyKeys(t *testing.T) {
	// BODY3D.ADDFORCE/ADDIMPULSE never existed; correct keys are APPLYFORCE/APPLYIMPULSE.
	forceAliases := []string{"addForce", "force", "applyForce"}
	for _, method := range forceAliases {
		k, p, o := handleCallBuiltin(heap.TagPhysicsBody, method)
		if !o || !p || k != "BODY3D.APPLYFORCE" {
			t.Fatalf("body3d %s: got %q prepend=%v ok=%v, want BODY3D.APPLYFORCE", method, k, p, o)
		}
	}
	impulseAliases := []string{"addImpulse", "impulse"}
	for _, method := range impulseAliases {
		k, p, o := handleCallBuiltin(heap.TagPhysicsBody, method)
		if !o || !p || k != "BODY3D.APPLYIMPULSE" {
			t.Fatalf("body3d %s: got %q prepend=%v ok=%v, want BODY3D.APPLYIMPULSE", method, k, p, o)
		}
	}
}

func TestHandleCallBuiltinTerrainSnappyMapsToSnapY(t *testing.T) {
	// TERRAIN.SNAPPY (two-P typo) never existed; correct key is TERRAIN.SNAPY.
	k, p, o := handleCallBuiltin(heap.TagTerrain, "snappy")
	if !o || !p || k != "TERRAIN.SNAPY" {
		t.Fatalf("terrain snappy: got %q prepend=%v ok=%v, want TERRAIN.SNAPY", k, p, o)
	}
}

func TestHandleCallBuiltinBody2DNewMethods(t *testing.T) {
	cases := []struct{ method, want string }{
		{"applyForce", "BODY2D.APPLYFORCE"},
		{"addForce", "BODY2D.APPLYFORCE"},
		{"applyImpulse", "BODY2D.APPLYIMPULSE"},
		{"addImpulse", "BODY2D.APPLYIMPULSE"},
		{"addCircle", "BODY2D.ADDCIRCLE"},
		{"addRect", "BODY2D.ADDRECT"},
		{"addPolygon", "BODY2D.ADDPOLYGON"},
		{"commit", "BODY2D.COMMIT"},
		{"collided", "BODY2D.COLLIDED"},
		{"collisionOther", "BODY2D.COLLISIONOTHER"},
		{"collisionPoint", "BODY2D.COLLISIONPOINT"},
		{"collisionNormal", "BODY2D.COLLISIONNORMAL"},
	}
	for _, tc := range cases {
		k, p, o := handleCallBuiltin(heap.TagBody2D, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("body2d %s: got %q prepend=%v ok=%v, want %q", tc.method, k, p, o, tc.want)
		}
	}
}

func TestHandleCallBuiltinBody3DNewMethods(t *testing.T) {
	cases := []struct{ method, want string }{
		{"activate", "BODY3D.ACTIVATE"},
		{"deactivate", "BODY3D.DEACTIVATE"},
		{"collided", "BODY3D.COLLIDED"},
		{"collisionOther", "BODY3D.COLLISIONOTHER"},
		{"collisionPoint", "BODY3D.COLLISIONPOINT"},
		{"collisionNormal", "BODY3D.COLLISIONNORMAL"},
		{"bufferIndex", "BODY3D.BUFFERINDEX"},
		{"setLinearVel", "BODY3D.SETLINEARVEL"},
		{"setLinearVelocity", "BODY3D.SETLINEARVEL"},
		{"getLinearVel", "BODY3D.GETLINEARVEL"},
		{"getLinearVelocity", "BODY3D.GETLINEARVEL"},
		{"applyImpulse", "BODY3D.APPLYIMPULSE"},
	}
	for _, tc := range cases {
		k, p, o := handleCallBuiltin(heap.TagPhysicsBody, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("body3d %s: got %q prepend=%v ok=%v, want %q", tc.method, k, p, o, tc.want)
		}
	}
}

func TestHandleCallBuiltinBBoxBSphereCollisionMethods(t *testing.T) {
	bboxCases := []struct{ method, want string }{
		{"check", "BBOX.CHECK"},
		{"checkSphere", "BBOX.CHECKSPHERE"},
	}
	for _, tc := range bboxCases {
		k, p, o := handleCallBuiltin(heap.TagBBox, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("bbox %s: got %q prepend=%v ok=%v, want %q", tc.method, k, p, o, tc.want)
		}
	}
	bsphereCases := []struct{ method, want string }{
		{"check", "BSPHERE.CHECK"},
		{"checkBox", "BSPHERE.CHECKBOX"},
	}
	for _, tc := range bsphereCases {
		k, p, o := handleCallBuiltin(heap.TagBSphere, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("bsphere %s: got %q prepend=%v ok=%v, want %q", tc.method, k, p, o, tc.want)
		}
	}
}

func TestHandleCallBuiltinCharControllerNewMethods(t *testing.T) {
	cases := []struct{ method, want string }{
		{"setLinearVelocity", "CHARACTERREF.SETLINEARVELOCITY"},
		{"setLinearVel", "CHARACTERREF.SETLINEARVELOCITY"},
		{"updateMove", "CHARACTERREF.UPDATEMOVE"},
		{"drainContacts", "CHARACTERREF.DRAINCONTACTS"},
		{"setContactListener", "CHARACTERREF.SETCONTACTLISTENER"},
		{"setSetting", "CHARACTERREF.SETSETTING"},
		{"getCeiling", "CHARACTERREF.GETCEILING"},
		{"getIsSliding", "CHARACTERREF.GETISSLIDING"},
		{"getGroundVelocity", "CHARACTERREF.GETGROUNDVELOCITY"},
	}
	for _, tc := range cases {
		k, p, o := handleCallBuiltin(heap.TagCharController, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("charcontroller %s: got %q prepend=%v ok=%v, want %q", tc.method, k, p, o, tc.want)
		}
	}
}

func TestHandleCallBuiltinWaterAutoPhysicsUpdate(t *testing.T) {
	cases := []struct{ method, want string }{
		{"autoPhysics", "WATER.AUTOPHYSICS"},
		{"update", "WATER.UPDATE"},
	}
	for _, tc := range cases {
		k, p, o := handleCallBuiltin(heap.TagWater, tc.method)
		if !o || !p || k != tc.want {
			t.Fatalf("water %s: got %q prepend=%v ok=%v, want %q", tc.method, k, p, o, tc.want)
		}
	}
}
