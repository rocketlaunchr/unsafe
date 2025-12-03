package unsafe

import (
	"reflect"
	"unsafe"
)

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

// Malloc allocates contiguous blocks of memory without zeroing the values.
// All allocated memory must immediately be filled by the developer, since
// the pre-existing data may contain garbage.
//
// Warning: If the returned pointer will be cast to a struct with fields that
// contain pointers, then provide a reflect Type for that struct.
// See: https://github.com/golang/go/issues/76352#issuecomment-3549768452
//
// Example:
//
//	type Person struct {
//	   name  string
//	   age   int
//	   phone *int
//	}
//
//	ptr := Malloc(unsafe.Sizeof(Person{}), false, reflect.TypeFor[Person]())
//	p := *(*Person)(ptr)
func Malloc(size uintptr, zero bool, t ...reflect.Type) unsafe.Pointer {
	if zero {
		return mallocgc(size, 0, true)
	}

	baseAdd := mallocgc(size, 0, false)

	if len(t) > 0 {
		for _, v := range findPointerFields(t[0], 0, true) {
			memclrNoHeapPointers(unsafe.Add(baseAdd, v.offset), v.size)
		}
	}
	return baseAdd
}

// New allocates memory for a struct.
// The value returned is a pointer to a newly allocated value of that type.
// The returned struct is not guaranteed to be the zero value.
//
// Example:
//
//	var pts = FindPointerFields(reflect.TypeFor[Person]())
//
//	type Person struct {
//	   name  string
//	   age   int
//	   phone *int
//	}
//
//	func main() {
//	   p := New[Person](pts) // like p := new(Person)
//	}
func New[T any](fp ...[]PtrOs) *T {
	var x *T
	var ptr unsafe.Pointer
	if len(fp) > 0 {
		ptr = Malloc(unsafe.Sizeof(*x), false)
		if fp[0] != nil {
			for _, v := range fp[0] {
				memclrNoHeapPointers(unsafe.Add(ptr, v.offset), v.size)
			}
		}
	} else {
		ptr = Malloc(unsafe.Sizeof(*x), false, reflect.TypeFor[T]())
	}
	return (*T)(ptr)
}

// NewZero allocates memory for a zero-value struct.
// The value returned is a pointer to a newly allocated value of that type.
// It is equivalent to new(T) (but appears to be faster).
func NewZero[T any]() *T {
	var x *T
	return (*T)(Malloc(unsafe.Sizeof(*x), true))
}
