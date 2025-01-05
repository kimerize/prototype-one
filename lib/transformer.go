package lib

import (
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Transformer interface {
	Transform(items []unstructured.Unstructured) []unstructured.Unstructured
}

type Generator interface {
	Generate() []unstructured.Unstructured
}

type ResourceGroup[T Transformer] struct {
	Config T
}

type dummyTransformer struct{}

func (dummyTransformer) Transform(items []unstructured.Unstructured) []unstructured.Unstructured {
	return items
}

var _ Generator = ResourceGroup[dummyTransformer]{}

func (t ResourceGroup[T]) Generate() []unstructured.Unstructured {
	v := reflect.ValueOf(t.Config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}

	items := []unstructured.Unstructured{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if g, ok := field.Interface().(Generator); ok {
			items = append(items, g.Generate()...)
		}
	}
	return t.Config.Transform(items)
}

func (t ResourceGroup[T]) WithOverrides(override func(*ResourceGroup[T])) ResourceGroup[T] {
	override(&t)
	return t
}

// func Generate[T any](fn func(T) []unstructured.Unstructured) Transform[T] {
// 	return func(items []unstructured.Unstructured, config T) []unstructured.Unstructured {
// 		return append(items, fn(config)...)
// 	}
// }

// type GeneratorList []Generator

// var _ Generator = GeneratorList{}

// func (l GeneratorList) Generate() []unstructured.Unstructured {
// 	items := []unstructured.Unstructured{}
// 	for _, g := range l {
// 		items = append(items, g.Generate()...)
// 	}
// 	return items
// }
