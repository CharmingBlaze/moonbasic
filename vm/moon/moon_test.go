package moon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"moonbasic/vm/opcode"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	p := opcode.NewProgram()
	p.Main.Emit(opcode.OpPushInt, p.Main.AddInt(7), 0, 1)
	p.Main.Emit(opcode.OpHalt, 0, 0, 2)

	fn := opcode.NewChunk("FOO")
	fn.Emit(opcode.OpReturnVoid, 0, 0, 10)
	p.Functions["FOO"] = fn

	p.Types["T"] = &opcode.TypeDef{Name: "T", Fields: []string{"X", "Y"}}

	data, err := Encode(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) < HeaderSize {
		t.Fatal("short file")
	}
	got, err := Decode(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Main.Instructions) != 2 {
		t.Fatalf("main instr: %d", len(got.Main.Instructions))
	}
	if got.Functions["FOO"] == nil {
		t.Fatal("missing FOO")
	}
	if got.Types["T"] == nil || len(got.Types["T"].Fields) != 2 {
		t.Fatal("missing typedef")
	}
}

func TestValidateHeaderRejectsGarbage(t *testing.T) {
	_, err := ValidateHeader(bytes.NewReader([]byte("XXXX\000\000\000\000\000\000\000\000\000\000\000\000")))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestValidateHeaderRejectsMoonV1(t *testing.T) {
	var hdr [HeaderSize]byte
	copy(hdr[0:4], []byte("MOON"))
	binary.BigEndian.PutUint32(hdr[4:8], VersionV1)
	binary.BigEndian.PutUint32(hdr[12:16], HeaderSize)
	_, err := ValidateHeader(bytes.NewReader(hdr[:]))
	if err == nil {
		t.Fatal("expected error for v1")
	}
	if !errors.Is(err, ErrVersion) {
		t.Fatalf("want ErrVersion wrap, got %v", err)
	}
}
