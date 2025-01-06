package lib

import (
	"reflect"

	"sigs.k8s.io/kustomize/api/resmap"
)

type Transformer interface {
	Transform(items resmap.ResMap)
}

type Generator interface {
	Generate() resmap.ResMap
}

type Overlay[T Transformer] struct {
	Config T
}

type dummyTransformer struct{}

func (dummyTransformer) Transform(items resmap.ResMap) {}

var _ Generator = Overlay[dummyTransformer]{}

func (t Overlay[T]) Generate() resmap.ResMap {
	v := reflect.ValueOf(t.Config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	result := resmap.New()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if g, ok := field.Interface().(Generator); ok {
			result.AbsorbAll(g.Generate())
		}
	}
	t.Config.Transform(result)
	return result
}

func (t Overlay[T]) WithOverrides(override func(*Overlay[T])) Overlay[T] {
	override(&t)
	return t
}

// func Generate[T any](fn func(T) resmap.ResMap) Transform[T] {
// 	return func(items resmap.ResMap, config T) resmap.ResMap {
// 		return append(items, fn(config)...)
// 	}
// }

// type GeneratorList []Generator

// var _ Generator = GeneratorList{}

// func (l GeneratorList) Generate() resmap.ResMap {
// 	items := resmap.ResMap{}
// 	for _, g := range l {
// 		items = append(items, g.Generate()...)
// 	}
// 	return items
// }
