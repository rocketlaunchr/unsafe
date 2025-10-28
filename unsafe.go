package unsafe

import (
	"reflect"
	"unsafe"
)

type Pointer = unsafe.Pointer
type IntegerType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type ArbitraryType interface {
	any
}

func Alignof[_ArbitraryType ArbitraryType](x _ArbitraryType) uintptr {
	return unsafe.Alignof(x)
}

func Sizeof[_ArbitraryType ArbitraryType](x _ArbitraryType) uintptr {
	return unsafe.Sizeof(x)
}

func String[_IntegerType IntegerType](ptr *byte, len _IntegerType) string {
	return unsafe.String(ptr, len)
}

func StringData(str string) *byte {
	return unsafe.StringData(str)
}

func Slice[_ArbitraryType ArbitraryType, _IntegerType IntegerType](ptr *_ArbitraryType, len _IntegerType) []_ArbitraryType {
	return unsafe.Slice(ptr, len)
}

func SliceData[_ArbitraryType ArbitraryType](slice []_ArbitraryType) *_ArbitraryType {
	return unsafe.SliceData(slice)
}

func Add[_IntegerType IntegerType](ptr Pointer, len _IntegerType) Pointer {
	return unsafe.Add(ptr, len)
}

// F calls the Field function directly.
func F[F interface{ int | []int | string }](f F, typ ...reflect.Type) field {
	return Field(f, typ...)
}

// Field selects which field in a struct you wish to modify.
// The field name can be selected by passing an integer or a string.
// An optional type constraint can be provided as a safety precaution to ensure that the field's type
// is what you expected.
func Field[F interface{ int | []int | string }](f F, typ ...reflect.Type) field {
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
	}
	panic("won't reach")
}

type field struct {
	i  *int
	is *[]int
	n  *string
	t  *reflect.Type
}

// SetField allows unexported fields of a struct to be modified.
// It returns a Pointer to the field if you wish to change the value
// externally instead of by passing a new value.
//
// NOTE: This function panics if strct is not actually a struct or the
// field could not be found.
func SetField[V any](strct any, f field, newValue ...V) Pointer {
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

	if f.i != nil {
		if f.t != nil && v.Field(*f.i).Type() != *f.t {
			panic("struct field type does not match")
		}
		ptr := v.Field(*f.i).Addr().UnsafePointer()
		if len(newValue) > 0 {
			*(*V)(ptr) = newValue[0]
		}
		return ptr
	} else if f.is != nil {
		if f.t != nil && v.Field(*f.i).Type() != *f.t {
			panic("struct field type does not match")
		}
		ptr := v.FieldByIndex(*f.is).Addr().UnsafePointer()
		if len(newValue) > 0 {
			*(*V)(ptr) = newValue[0]
		}
		return ptr
	} else if f.n != nil {
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Name == *f.n {
				if f.t != nil && v.Field(i).Type() != *f.t {
					panic("struct field type does not match")
				}
				ptr := v.Field(i).Addr().UnsafePointer()
				if len(newValue) > 0 {
					*(*V)(ptr) = newValue[0]
				}
				return ptr
			}
		}
		panic("struct field name not found")
	}
	panic("won't reach")
}
