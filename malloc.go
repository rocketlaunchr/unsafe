package unsafe

import (
	"reflect"
	"unsafe"
)

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ uintptr, needzero bool) unsafe.Pointer

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
		return unsafe.Pointer(unsafe.SliceData(make([]byte, size)))
	}

	baseAdd := mallocgc(size, 0, false)

	if len(t) > 0 {
		for _, v := range findPointerFields(t[0], 0, true) {
			SetZeros(baseAdd, v.offset, v.size)
		}
	}
	return baseAdd
}

// New allocates memory for a struct.
// The value returned is a pointer to a newly allocated value of that type.
// The returned struct is not guaranteed to be the zero value.
func New[T any](fp ...[]pointers) *T {
	var x *T
	var ptr unsafe.Pointer
	if len(fp) > 0 {
		ptr = Malloc(unsafe.Sizeof(*x), false)
		if fp[0] != nil {
			for _, v := range fp[0] {
				SetZeros(ptr, v.offset, v.size)
			}
		}
	} else {
		ptr = Malloc(unsafe.Sizeof(*x), false, reflect.TypeFor[T]())
	}
	return (*T)(ptr)
}

// SetZeros sets memory address located at (baseAddress + offset) to zero.
// Size (in bytes) dictates how many contiguous locations to set to zero.
// No bytes are set to zero if size is 0.
func SetZeros(baseAddress unsafe.Pointer, offset, size uintptr) {
	address := unsafe.Add(baseAddress, offset)
	switch size {
	case 0:
		break
	case 1:
		*(*uint8)(address) = 0
	case 2:
		*(*uint16)(address) = 0
	case 4:
		*(*uint32)(address) = 0
	case 32, 40, 48, 56, 64:
		mult := size / 8
		for i := uintptr(0); i < mult-1; i++ {
			*(*uint64)(unsafe.Add(address, 8*(i+1))) = 0
		}
		fallthrough
	case 8:
		*(*uint64)(address) = 0
	case 16:
		*(*uint64)(address) = 0
		*(*uint64)(unsafe.Add(address, 8)) = 0
	case 24:
		*(*uint64)(address) = 0
		*(*uint64)(unsafe.Add(address, 8)) = 0
		*(*uint64)(unsafe.Add(address, 16)) = 0
	default:
		for i := uintptr(0); i < size; i++ {
			*(*byte)(unsafe.Add(address, i)) = 0
		}
	}
}
