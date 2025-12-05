package unsafe

import (
	"unsafe"
)

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

// Malloc allocates contiguous blocks of memory without zeroing.
// However, the memory will be cleared when zerofill is set.
//
// All allocated memory must sensibly be filled by the developer immediately,
// since the pre-existing data may contain garbage.
//
// WARNING: If the returned pointer will be cast to a struct with fields
// containing pointers, then data stored in those specific
// fields must be cleared before casting.
//
// See: https://github.com/golang/go/issues/76352#issuecomment-3549768452
//
// Note: Reference types (eg. strings, slices, interfaces, channels, maps etc.)
// contain pointers.
//
// Example:
//
//	type Person struct {
//	   name  string
//	   age   int
//	   phone *int
//	}
//
//	ptr := Malloc(unsafe.Sizeof(Person{}), false, FindPointerFields(reflect.TypeFor[Person]()))
//	p := *(*Person)(ptr)
func Malloc(size uintptr, zerofill bool, zero ...[][2]uintptr) unsafe.Pointer {
	base := mallocgc(size, 0, zerofill)
	if !zerofill && len(zero) > 0 {
		for _, v := range zero[0] {
			memclrNoHeapPointers(unsafe.Add(base, v[0]), v[1])
		}
	}
	return base
}

// UnsafePointerFields returns the offset and size for each field of a struct
// that contain pointers.
//
// Note: Reference types (eg. strings, slices, interfaces, channels, maps etc.) contain pointers.
type UnsafePointerFields interface {
	// UnsafePointerFields must be implemented on the pointer type *T.
	// The return values can be hardcoded using unsafe.Offsetof() and unsafe.Sizeof().
	// Alternatively, the automatic FindPointerFields function can be stored in a global
	// variable and subsequently returned.
	//
	// NOTE: When hardcoding, special attention is needed for fields of type Array with elements
	// containing pointers (including structs which itself contains pointer fields).
	UnsafePointerFields() [][2]uintptr
}

// New allocates memory for a struct without zeroing.
// A pointer to the newly allocated value is returned.
// The returned struct is not guaranteed to be the zero value.
//
// NOTE: UnsafePointerFields must be implemented as a method with a pointer receiver.
// UnsafePointerFields is used to selectively zero fields with pointers or reference types.
//
// Example:
//
//	var pfs = FindPointerFields(reflect.TypeFor[Person]())
//
//	type Person struct {
//	   name  string
//	   age   int
//	   phone *int
//	}
//
//	func (p *Person) UnsafePointerFields() [][2]uintptr {
//		return pfs
//		return [][2]uintptr {
//			{unsafe.Offsetof(Person{}.name), unsafe.Sizeof(Person{}.name)},
//			{unsafe.Offsetof(Person{}.phone), unsafe.Sizeof(Person{}.phone)},
//		}
//		return nil
//	}
//
//	func main() {
//	   p := New[Person]() // like p := new(Person)
//	}
func New[T any, PT interface {
	*T
	UnsafePointerFields
}]() *T {
	var x PT
	ptr := mallocgc(unsafe.Sizeof(*x), 0, false)
	for _, v := range (x).UnsafePointerFields() {
		memclrNoHeapPointers(unsafe.Add(ptr, v[0]), v[1])
	}

	return (*T)(ptr)
}

// NewZero allocates memory for a zero-value struct.
// A pointer to the newly allocated value is returned.
// It is equivalent to new(T) (but seems to be faster).
func NewZero[T any]() *T {
	var x *T
	return (*T)(mallocgc(unsafe.Sizeof(*x), 0, true))
}
