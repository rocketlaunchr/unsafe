// go test -count=1 -v -bench=. -benchmem

package unsafe

import (
	"reflect"
	"testing"
	"unsafe"
)

type Person struct {
	name  string
	age   int
	phone *int
}

var pfs = FindPointerFields(reflect.TypeFor[Person]())

func (p *Person) UnsafePointerFields() [][2]uintptr {
	return pfs
	return [][2]uintptr{
		{unsafe.Offsetof(Person{}.name), unsafe.Sizeof(Person{}.name)},
		{unsafe.Offsetof(Person{}.phone), unsafe.Sizeof(Person{}.phone)},
	}
	return nil
}

var result *Person

const N = 1000

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
			r = New[Person]()
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

func BenchmarkNewZero(b *testing.B) {
	var r *Person
	for n := 0; n < b.N; n++ {
		for n := 0; n < N; n++ {
			r = NewZero[Person]()
		}
	}
	result = r
}
