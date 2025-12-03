// go test -count=1 -v -bench=.

package unsafe

import (
	"reflect"
	"sync"
	"testing"
	"unsafe"
)

type Person struct {
	name  string
	age   int
	phone *int
}

var result *Person

const N = 1000

var pts = sync.OnceValue(func() []PtrOs {
	return FindPointerFields(reflect.TypeFor[Person]())
})

var pts2 = FindPointerFields(reflect.TypeFor[Person]())

func init() {
	pts()
}

func BenchmarkMallocNew(b *testing.B) {
	var r *Person
	for n := 0; n < b.N; n++ {
		for n := 0; n < N; n++ {
			ptr := Malloc(unsafe.Sizeof(Person{}), false)
			r = (*Person)(ptr)
		}
	}
	result = r
}

func BenchmarkMallocNewSelectiveZeroing(b *testing.B) {
	var r *Person
	for n := 0; n < b.N; n++ {
		for n := 0; n < N; n++ {
			// 			r = New[Person](pts()) // sync.Once is super-slow
			r = New[Person](pts2)
			// 			r = New[Person](nil) // Don't perform zeroing
		}
	}
	result = r
}

func BenchmarkStdNew(b *testing.B) {
	var r *Person
	for n := 0; n < b.N; n++ {
		for n := 0; n < N; n++ {
			r = new(Person)
		}
	}
	result = r
}
