package unsafe

import (
	"reflect"
	"unsafe"
)

type (
	// Type
	//
	// See: https://pkg.go.dev/reflect#Type
	Type = reflect.Type

	// Pointer
	//
	// See: https://pkg.go.dev/unsafe#Pointer
	Pointer = unsafe.Pointer

	_IntegerType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}

	_ArbitraryType interface {
		any
	}
)

// Alignof
//
// See: https://pkg.go.dev/unsafe#Alignof
func Alignof[ArbitraryType _ArbitraryType](x ArbitraryType) uintptr {
	return unsafe.Alignof(x)
}

// Sizeof
//
// See: https://pkg.go.dev/unsafe#Sizeof
func Sizeof[ArbitraryType _ArbitraryType](x ArbitraryType) uintptr {
	return unsafe.Sizeof(x)
}

// String
//
// See: https://pkg.go.dev/unsafe#String
func String[IntegerType _IntegerType](ptr *byte, len IntegerType) string {
	return unsafe.String(ptr, len)
}

// StringData
//
// See: https://pkg.go.dev/unsafe#StringData
func StringData(str string) *byte {
	return unsafe.StringData(str)
}

// Slice
//
// See: https://pkg.go.dev/unsafe#Slice
func Slice[ArbitraryType _ArbitraryType, IntegerType _IntegerType](ptr *ArbitraryType, len IntegerType) []ArbitraryType {
	return unsafe.Slice(ptr, len)
}

// SliceData
//
// See: https://pkg.go.dev/unsafe#SliceData
func SliceData[ArbitraryType _ArbitraryType](slice []ArbitraryType) *ArbitraryType {
	return unsafe.SliceData(slice)
}

// Add
//
// See: https://pkg.go.dev/unsafe#Add
func Add[IntegerType _IntegerType](ptr Pointer, len IntegerType) Pointer {
	return unsafe.Add(ptr, len)
}

// TypeOf returns the reflection [Type] that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i any) Type {
	return reflect.TypeOf(i)
}

// TypeFor returns the [Type] that represents the type argument T.
func TypeFor[T any]() Type {
	return reflect.TypeFor[T]()
}

type field struct {
	i  *int
	is *[]int
	n  *string
	u  *uintptr
	t  *Type
}

// F is shorthand for the Field function.
func F[F interface{ int | []int | string | uintptr }](f F, typ ...Type) field {
	return Field(f, typ...)
}

// Field selects a field in a struct.
// The field can be selected by passing an integer or a (name) string.
// An optional type constraint can be provided as a safety precaution to ensure that the field's type
// is what you expected.
func Field[F interface{ int | []int | string | uintptr }](f F, typ ...Type) field {
	var t *reflect.Type
	if len(typ) > 0 {
		t = &(typ[0])
	}

	switch f := any(f).(type) {
	case int:
		return field{i: &f, t: t}
	case []int:
		return field{is: &f, t: t}
	case string:
		return field{n: &f, t: t}
	case uintptr:
		return field{u: &f, t: t}
	}
	panic("won't reach")
}

// Value returns the value of the unexported field of a struct.
//
// NOTE: This function panics if strct is not actually a struct or the
// field could not be found.
func Value[T any](strct any, f field) T {
	ptr := SetField[T](strct, f)
	return *(*T)(ptr)
}

// SetField allows unexported fields of a struct to be modified.
// It returns a Pointer to the field if you wish to change the value
// externally instead of by passing a new value.
//
// NOTE: This function panics if strct is not actually a struct or the
// field could not be found.
func SetField[T any](strct any, f field, newValue ...T) Pointer {
	// Special case - strct is an unsafe.Pointer
	if ptr, ok := strct.(unsafe.Pointer); ok {
		if f.u != nil {
			ptr = unsafe.Add(ptr, *f.u)
		}

		if len(newValue) > 0 {
			*(*T)(ptr) = newValue[0]
		}
		return ptr
	}

	v := reflect.ValueOf(strct)
	if !v.IsValid() {
		panic("strct not valid")
	}
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			panic("strct is nil pointer")
		}
		v = v.Elem()
	}
	t := v.Type()
	if t.Kind() != reflect.Struct {
		panic("strct is not a struct")
	}

	var ptr unsafe.Pointer

	if f.i != nil {
		if f.t != nil && v.Field(*f.i).Type() != *f.t {
			panic("struct field type does not match")
		}
		ptr = v.Field(*f.i).Addr().UnsafePointer()
	} else if f.is != nil {
		if f.t != nil && v.FieldByIndex(*f.is).Type() != *f.t {
			panic("struct field type does not match")
		}
		ptr = v.FieldByIndex(*f.is).Addr().UnsafePointer()
	} else if f.n != nil {
		fbn := v.FieldByName(*f.n)
		if fbn.IsValid() {
			if f.t != nil && fbn.Type() != *f.t {
				panic("struct field type does not match")
			}
			ptr = fbn.Addr().UnsafePointer()
		} else {
			panic("struct field name not found")
		}
	} else if f.u != nil {
		ptr = unsafe.Add(v.Addr().UnsafePointer(), *f.u)
	}
	if len(newValue) > 0 {
		*(*T)(ptr) = newValue[0]
	}
	return ptr
}
