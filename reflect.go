package unsafe

import (
	"reflect"
	"slices"
)

// FindPointerFields finds in a struct the offsets (and corresponding sizes)
// of all fields that contain pointers and other reference types.
func FindPointerFields(t reflect.Type) [][2]uintptr {
	return findPointerFields(t, 0, true)
}

func findPointerFields(t reflect.Type, currentOffset uintptr, top bool) [][2]uintptr {
	var results [][2]uintptr
	defer func() { merge(&results) }()

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
	NextField:
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
				if maxElements == 0 {
					continue NextField
				}
				f = f.Elem()

				switch f.Kind() {
				case reflect.Array:
				// Not possible
				case reflect.Slice, reflect.String, reflect.Interface, reflect.Chan, reflect.Map, reflect.Func, reflect.Pointer, reflect.UnsafePointer:
					results = append(results, [2]uintptr{fieldOffset, field.Type.Size()})
				case reflect.Struct:
					pfs := findPointerFields(f, fieldOffset, false)
					res := make([][2]uintptr, 0, len(pfs)*maxElements)
					for _, val := range pfs {
						for i := 0; i < maxElements; i++ {
							res = append(res, [2]uintptr{val[0] + f.Size()*uintptr(i), val[1]})
						}
					}
					results = append(results, res...)
				default:
					// Non-reference types
				}
			case reflect.Slice, reflect.String, reflect.Interface, reflect.Chan, reflect.Map, reflect.Func, reflect.Pointer, reflect.UnsafePointer:
				results = append(results, [2]uintptr{fieldOffset, field.Type.Size()})
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

func merge(s *[][2]uintptr) {
	toDel := []int{}
	for i := len(*s) - 1; i > 0; i-- {
		current := &(*s)[i]
		prev := &(*s)[i-1]

		if prev[0]+prev[1] == current[0] {
			prev[1] += current[1]
			toDel = append(toDel, i)
		}
	}

	for _, val := range toDel {
		*s = slices.Delete(*s, val, val+1)
	}
}
