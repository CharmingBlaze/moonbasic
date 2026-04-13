//go:build cgo || (windows && !cgo)

package rl

import (
	"fmt"
	"reflect"
	"unsafe"
)

// VertexAttributesConfig is used by [SetVertexAttributes] to specify VAO bindings for a slice of structs or arrays.
type VertexAttributesConfig struct {
	Field      string // Name of the field in the struct (ignored when slice is backed by an array instead of struct [][2]float32)
	Attribute  uint32 // OpenGL attribute index (layout location)
	Normalized bool   // Whether the attribute should be normalized
}

// SetVertexAttributes can automatically define VAO bindings for a slice of structs or slice of 1d arrays, with supported primitive types.
// Supported primitives: float32, float64m and signed/unsigned integers (except int, int64),
// NOTE: bind VAO and VBO before calling this.
//
// If a slice of structs is passed []VertexStruct:
// The struct can contain primitives, or an array of primtives, or a struct that's
// made out of the same primitve for every field.
//
//	type VertexStruct struct{
//		Position rl.Vector3 // has same primitive type for all fields (float32), hence allowed.
//		Color [4]byte       // array of 4 bytes. Also allowed.
//		BlockId uint8       // a primitive. Allowed
//	}
//	var Vertices := [][4]float32{} // allowed. Slice of array of primitives.
func SetVertexAttributes[T any](vertices []T, attributes []VertexAttributesConfig) {
	if len(vertices) == 0 {
		return
	}
	// reflect.TypeFor but for go 1.21
	t := reflect.TypeOf((*T)(nil)).Elem()
	// Compute stride (size of one vertex in bytes)
	stride := int32(unsafe.Sizeof(vertices[0]))
	kind := t.Kind()

	switch kind {
	default:
		panic("Vertex array is using unsupported types. Only structs and and arrays are supported.")
	case reflect.Array: // slice of arrays eg. [][2]float32
		arrayKind := t.Elem().Kind()               // backing type of the array. eg. float32
		attrType, isPrimitive := glType(arrayKind) // convert to GL type.

		if !isPrimitive { // type could not be converted because unsupported by GL
			panic("Backing type for array is not one of the supported primitives " + t.Elem().String())
		}
		components := int32(t.Len())     // each attribute has same number of components (len of array)
		attributeSize := t.Elem().Size() //each attribute has same size.

		// iterate over each vertex attribute.
		for i, attr := range attributes {
			offset := int32(i) * int32(attributeSize) // manually calculate offset
			// call OpenGL to define this vertex attribute
			SetVertexAttribute(attr.Attribute, components, attrType, attr.Normalized, stride, offset)
			EnableVertexAttribute(attr.Attribute)
		}
	//* A struct can contain:
	// a primtive
	// an array of primitives
	// a struct that's made out of the same primtives for every field (basically a named array)
	case reflect.Struct:
		// Iterate over each attribute configuration
		for _, attr := range attributes {
			// Find the field by name
			field, ok := t.FieldByName(attr.Field)
			if !ok {
				panic(fmt.Sprintf("struct %s does not have a field of the name %s", t.String(), attr.Field))
			}

			// Check if the field is a primitive type (float32, uint8, etc.)
			attrType, isPrimitiveType := glType(field.Type.Kind())
			if isPrimitiveType {
				components := int32(1)
				offset := int32(field.Offset)

				// call OpenGL to define this vertex attribute
				SetVertexAttribute(attr.Attribute, components, attrType, attr.Normalized, stride, offset)
				EnableVertexAttribute(attr.Attribute)
				continue
			}
			// Field is not a primitive. Check if the field is an array of primitives.
			switch field.Type.Kind() {
			case reflect.Array:
				// Array of primitive types
				elemKind := field.Type.Elem().Kind()
				components := int32(field.Type.Len())
				offset := int32(field.Offset)
				attrType, isPrimitiveType := glType(elemKind) // check if array of primitives.
				if !isPrimitiveType {
					panic(fmt.Sprint("Only array of primitive types is supported. Got ", elemKind.String(), " for field ", attr.Field))
				}
				// call OpenGL
				SetVertexAttribute(attr.Attribute, components, attrType, attr.Normalized, stride, offset)
				EnableVertexAttribute(attr.Attribute)
			// field is not an array of primitives. Is it a struct instead?
			// Each field in this child struct must be of the same primitive type.
			// The child struct is basically treated like an array.
			case reflect.Struct:
				components := int32(field.Type.NumField()) // how many primitives?
				if components == 0 {
					panic(fmt.Sprintf("Child struct %s is empty in field %s", field.Type.String(), attr.Field))
				}
				// Child struct: ensure all fields are the same primitive type
				prevType := field.Type.Field(0).Type
				attrType, isPrimitiveType := glType(prevType.Kind())
				offset := int32(field.Offset) // offset of this field within the vertex

				if !isPrimitiveType {
					panic(fmt.Sprintf("child struct must have a primitive type for every field in %s %s", attr.Field, prevType.String()))
				}

				// Check that all child struct fields have the same type
				for i := 1; i < field.Type.NumField(); i++ {
					if prevType != field.Type.Field(i).Type {
						panic(fmt.Sprintf("child struct must have the same type for every field in %s %s", attr.Field, field.Type.String()))
					}
				}

				// call OpenGL
				SetVertexAttribute(attr.Attribute, components, attrType, attr.Normalized, stride, offset)
				EnableVertexAttribute(attr.Attribute)
			}
		}
	}
}

func glType(k reflect.Kind) (t int32, ok bool) {
	switch k {
	case reflect.Int8:
		return Byte, true
	case reflect.Uint8:
		return UnsignedByte, true
	case reflect.Int16:
		return Short, true
	case reflect.Uint16:
		return UnsignedShort, true
	case reflect.Int32:
		return Int, true
	case reflect.Uint32:
		return UnsignedInt, true
	case reflect.Float32:
		return Float, true
	case reflect.Float64:
		return Double, true
	default:
		return -1, false
	}
}
