package mbdata

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"

	"github.com/klauspost/compress/zstd"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("DATA.COMPRESS", "data", m.dataCompress)
	r.Register("DATA.DECOMPRESS", "data", m.dataDecompress)
	r.Register("DATA.ENCODEBASE64", "data", m.dataEncodeBase64)
	r.Register("DATA.DECODEBASE64", "data", m.dataDecodeBase64)
	r.Register("DATA.CRC32", "data", m.dataCRC32)
	r.Register("DATA.MD5", "data", m.dataMD5)
	r.Register("DATA.SHA1", "data", m.dataSHA1)
	r.Register("DATA.COMPUTECRC32", "data", m.dataComputeCRC32)
	r.Register("DATA.COMPUTEMD5", "data", m.dataComputeMD5)
	r.Register("DATA.COMPUTESHA1", "data", m.dataComputeSHA1)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) dataCompress(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.COMPRESS expects string argument")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	src := []byte(s)
	enc, err := zstd.NewWriter(nil)
	if err != nil {
		return value.Nil, fmt.Errorf("DATA.COMPRESS: %w", err)
	}
	defer enc.Close()
	out := enc.EncodeAll(src, make([]byte, 0, len(src)))
	return rt.RetString(string(out)), nil
}

func (m *Module) dataDecompress(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.DECOMPRESS expects string argument")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	src := []byte(s)
	dec, err := zstd.NewReader(nil)
	if err != nil {
		return value.Nil, fmt.Errorf("DATA.DECOMPRESS: %w", err)
	}
	defer dec.Close()
	out, err := dec.DecodeAll(src, nil)
	if err != nil {
		return value.Nil, fmt.Errorf("DATA.DECOMPRESS: %w", err)
	}
	return rt.RetString(string(out)), nil
}

func (m *Module) dataEncodeBase64(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.ENCODEBASE64 expects string argument")
	}
	srcStr, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	src := []byte(srcStr)
	encoded := base64.StdEncoding.EncodeToString(src)
	return rt.RetString(encoded), nil
}

func (m *Module) dataDecodeBase64(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.DECODEBASE64 expects string argument")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	out, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return value.Nil, fmt.Errorf("DATA.DECODEBASE64: %w", err)
	}
	return rt.RetString(string(out)), nil
}

func (m *Module) dataCRC32(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.CRC32 expects string argument")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	src := []byte(s)
	c := crc32.ChecksumIEEE(src)
	return value.FromInt(int64(uint32(c))), nil
}

func (m *Module) dataMD5(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.MD5 expects string argument")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	src := []byte(s)
	sum := md5.Sum(src)
	return rt.RetString(hex.EncodeToString(sum[:])), nil
}

func (m *Module) dataSHA1(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("DATA.SHA1 expects string argument")
	}
	s, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	src := []byte(s)
	sum := sha1.Sum(src)
	return rt.RetString(hex.EncodeToString(sum[:])), nil
}

func (m *Module) dataComputeCRC32(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.dataCRC32(rt, args...)
}

func (m *Module) dataComputeMD5(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.dataMD5(rt, args...)
}

func (m *Module) dataComputeSHA1(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.dataSHA1(rt, args...)
}
