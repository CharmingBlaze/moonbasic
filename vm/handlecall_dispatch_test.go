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
}

func TestHandleCallBuiltinBody2DRotSetter(t *testing.T) {
	k, prep, ok := handleCallBuiltin(heap.TagBody2D, "rot")
	if !ok || !prep || k != "BODY2D.SETROT" {
		t.Fatalf("body2d rot(angle): got %q prepend=%v ok=%v", k, prep, ok)
	}
}
