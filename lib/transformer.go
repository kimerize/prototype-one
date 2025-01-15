package lib

import (
	"reflect"
)

type Transformer interface {
	Transform(items *ResourceList)
}

type Generator interface {
	Generate() ResourceList
}

type Overlay[T Transformer] struct {
	Config T
}

type dummyTransformer struct {
	kind string
}

func (dummyTransformer) Transform(items *ResourceList) {}

var _ Generator = Overlay[dummyTransformer]{}

func (t Overlay[T]) Generate() ResourceList {
	v := reflect.ValueOf(t.Config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("Overlay config must be a struct")
	}

	result := NewResourceList()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if g, ok := field.Interface().(Generator); ok {
			result.Absorb(g.Generate())
		}
	}
	t.Config.Transform(result)
	return *result
}

func (t Overlay[T]) WithOverrides(override func(*Overlay[T])) Overlay[T] {
	override(&t)
	return t
}

// func Generate[T any](fn func(T) ResourceList) Transform[T] {
// 	return func(items ResourceList, config T) ResourceList {
// 		return append(items, fn(config)...)
// 	}
// }

// type GeneratorList []Generator

// var _ Generator = GeneratorList{}

// func (l GeneratorList) Generate() ResourceList {
// 	items := ResourceList{}
// 	for _, g := range l {
// 		items = append(items, g.Generate()...)
// 	}
// 	return items
// }
