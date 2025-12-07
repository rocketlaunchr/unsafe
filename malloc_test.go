// go test -count=1 -v -bench=. -benchmem

package unsafe

import (
	"testing"
	"unsafe"
)

type Person struct {
	name  [50]byte
	age   int
	phone int
}

var clrout = [][2]uintptr{
	{unsafe.Offsetof(Person{}.name), unsafe.Sizeof(Person{}.name)},
	{unsafe.Offsetof(Person{}.phone), unsafe.Sizeof(Person{}.phone)},
}

func (p *Person) UnsafeClearOutFields() [][2]uintptr {
	return clrout
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
