package unsafe

import (
	"reflect"
	"unsafe"
)

type goTyp struct {
	Size       uintptr
	PtrData    uintptr
	Hash       uint32
	Flags      uint8
	Align      uint8
	FieldAlign uint8
	KindFlags  uint8
	Traits     unsafe.Pointer
	GCData     *byte
	Str        int32
	PtrToSelf  int32
}

type goItab struct {
	it unsafe.Pointer
	Vt *goTyp
	hv uint32
	_  [4]byte
	fn [1]uintptr
}

type goIface struct {
	Itab  *goItab
	Value unsafe.Pointer
}

func goType(t reflect.Type) *goTyp {
	return (*goTyp)((*goIface)(unsafe.Pointer(&t)).Value)
}

//go:linkname mallocgc runtime.mallocgc
func mallocgc(size uintptr, typ *goTyp, needzero bool) unsafe.Pointer

//go:linkname memclrNoHeapPointers runtime.memclrNoHeapPointers
func memclrNoHeapPointers(ptr unsafe.Pointer, n uintptr)

// ClearOut zeros out selected blocks of memory.
// clearout is a slice containing offset-size pairs.
// The offset is always relative to baseAddr.
//
// NOTE: It uses runtime.memclrNoHeapPointers so make sure you read
// documentation on when it is safe to use.
func ClearOut(baseAddr unsafe.Pointer, clearout [][2]uintptr) {
	for _, v := range clearout {
		memclrNoHeapPointers(unsafe.Add(baseAddr, v[0]), v[1])
	}
}

// Malloc allocates contiguous blocks of memory without zeroing.
// However, the memory will be cleared when zerofill is set.
//
// All allocated memory must sensibly be filled by the developer immediately,
// since the pre-existing data may contain garbage.
//
// WARNING: If the returned pointer will be cast to a struct with fields containing
// pointers, then typ must be provided to inform the garbage collector.
// As an optimization, typ can be stored in a global variable and reused.
//
// NOTE: Reference types (eg. strings, slices, interfaces, channels, maps etc.)
// contain pointers.
//
// Example:
//
//	type Person struct {
//	   name  [50]byte
//	   age   int
//	   phone int
//	}
//
//	ptr := Malloc(unsafe.Sizeof(Person{}), false)
//	ClearOut(ptr, [][2]uintptr{})
//	p := *(*Person)(ptr)
func Malloc(size uintptr, zerofill bool, typ ...reflect.Type) unsafe.Pointer {
	var t *goTyp
	if len(typ) > 0 {
		t = goType(typ[0])
	}
	return mallocgc(size, t, zerofill)
}

// UnsafeClearOutFields is used to return the offset and size of each field of a struct
// that should be cleared.
type UnsafeClearOutFields interface {
	// UnsafeClearOutFields must be implemented on the pointer type *T.
	// The return values can be hardcoded using unsafe.Offsetof() and unsafe.Sizeof().
	UnsafeClearOutFields() [][2]uintptr
}

// New allocates memory for a struct without zeroing.
// A pointer to the newly allocated value is returned.
// The returned struct is not guaranteed to be the zero value.
//
// NOTE: UnsafeClearOutFields must be implemented as a method with a pointer receiver.
// UnsafeClearOutFields is used to selectively zero fields.
//
// WARNING: T must not have fields containing pointers. ...But if it
// does contain pointers, provide typ=nil to inform the garbage collector.
// As an optimization, typ can be stored in a global variable and reused.
//
// NOTE: Reference types (eg. strings, slices, interfaces, channels, maps etc.)
// contain pointers.
//
// Example:
//
//	type Person struct {
//	   name  [50]byte
//	   age   int
//	   phone int
//	}
//
//	func (p *Person) UnsafeClearOutFields() [][2]uintptr {
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
func New[T any, PtrT interface {
	*T
	UnsafeClearOutFields
}](typ ...reflect.Type,
) *T {
	var x PtrT
	var ptr unsafe.Pointer
	if len(typ) > 0 {
		if typ[0] == nil {
			typ[0] = reflect.TypeFor[T]()
		}
		ptr = mallocgc(unsafe.Sizeof(*x), goType(typ[0]), false)
	} else {
		ptr = mallocgc(unsafe.Sizeof(*x), nil, false)
	}

	for _, v := range x.UnsafeClearOutFields() {
		memclrNoHeapPointers(unsafe.Add(ptr, v[0]), v[1])
	}

	return (*T)(ptr)
}

// NewZero allocates memory for a zero-value struct.
// A pointer to the newly allocated value is returned.
// It is equivalent to new(T) (but seems to be faster).
//
// WARNING: T must not have fields containing pointers. ...But if it
// does contain pointers, provide typ=nil to inform the garbage collector.
// As an optimization, typ can be stored in a global variable and reused.
//
// NOTE: Reference types (eg. strings, slices, interfaces, channels, maps etc.)
// contain pointers.
//
// Example:
//
//	p := NewZero[T]()
//
// Deprecated: From Go1.26+, new(T) will be on par or faster than NewZero.
func NewZero[T any](typ ...reflect.Type) *T {
	var x *T
	if len(typ) > 0 {
		if typ[0] == nil {
			typ[0] = reflect.TypeFor[T]()
		}
		return (*T)(mallocgc(unsafe.Sizeof(*x), goType(typ[0]), true))
	}
	return (*T)(mallocgc(unsafe.Sizeof(*x), nil, true))
}
