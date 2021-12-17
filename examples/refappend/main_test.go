package main

import (
	"reflect"
	"testing"
)

func BenchmarkByAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 5)
		v1 := reflect.ValueOf(&s).Elem()
		v1 = reflect.Append(v1, reflect.ValueOf(1))
		v1 = reflect.Append(v1, reflect.ValueOf(2))
	}
}

func BenchmarkByAppend2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 5)
		v1 := reflect.ValueOf(&s).Elem()
		v1.Set(reflect.Append(v1, reflect.ValueOf(1)))
		v1.Set(reflect.Append(v1, reflect.ValueOf(2)))
	}
}

func BenchmarkBySet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := make([]int, 0, 5)
		v1 := reflect.ValueOf(&s).Elem()
		v1.SetLen(2)
		v1.Index(0).Set(reflect.ValueOf(1))
		v1.Index(1).Set(reflect.ValueOf(2))
	}
}
