package unsafe

import (
	"reflect"
)

type pointers struct {
	offset uintptr
	size   uintptr
}

// FindPointerFields finds the offsets of all fields in a struct
// that contain pointers and other reference types.
func FindPointerFields(t reflect.Type) []pointers {
	return findPointerFields(t, 0, true)
}

func findPointerFields(t reflect.Type, currentOffset uintptr, top bool) []pointers {
	var results []pointers

	if top {
		// If the type is a pointer, dereference it.
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		// Only process structs.
		if t.Kind() != reflect.Struct {
			return results
		}
	}

	fieldOffset := uintptr(0)

	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fieldOffset = currentOffset + field.Offset

			switch field.Type.Kind() {
			case reflect.Array:
				maxElements := 1
				f := field.Type
				for ; f.Elem().Kind() == reflect.Array; f = f.Elem() {
					maxElements = maxElements * f.Len()
				}
				maxElements = maxElements * f.Len()
				f = f.Elem()

				switch f.Kind() {
				case reflect.Array:
				// Not possible
				case reflect.Slice, reflect.String, reflect.Interface, reflect.Chan, reflect.Map, reflect.Func, reflect.Pointer, reflect.UnsafePointer:
					results = append(results, pointers{
						offset: (fieldOffset),
						size:   (field.Type.Size()),
					})
				case reflect.Struct:
					if maxElements > 0 {
						pfs := findPointerFields(f, fieldOffset, false)
						res := make([]pointers, 0, len(pfs)*maxElements)
						for _, val := range pfs {
							for i := 0; i < maxElements; i++ {
								res = append(res, pointers{
									offset: val.offset + f.Size()*uintptr(i),
									size:   val.size,
								})
							}
						}
						results = append(results, res...)
					}
				default:

				}
			case reflect.Slice, reflect.String, reflect.Interface, reflect.Chan, reflect.Map, reflect.Func, reflect.Pointer, reflect.UnsafePointer:
				results = append(results, pointers{
					offset: (fieldOffset),
					size:   (field.Type.Size()),
				})
			case reflect.Struct:
				results = append(results, findPointerFields(field.Type, fieldOffset, false)...)
			default:
				// We don't zero out these types since they aren't pointers:
				// Bool, Int, Int8, Int16, Int32, Int64, Uint, Uint8, Uint16,
				// Uint32, Uint64, Uintptr, Float32, Float64, Complex64, Complex128
			}
		}
	default:
		panic("encountered unexpected type: " + t.Kind().String())
	}

	return results
}
