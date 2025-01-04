package lib

import (
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Transform[T any] func(items []unstructured.Unstructured, config T) []unstructured.Unstructured

type Generator interface {
	Generate() []unstructured.Unstructured
}

type ResourceGroup[T any] struct {
	Transform[T]
	Config T
}

var _ Generator = ResourceGroup[any]{}

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
	return t.Transform(items, t.Config)
}

func (t ResourceGroup[T]) WithOverrides(override func(*ResourceGroup[T])) ResourceGroup[T] {
	override(&t)
	return t
}

func Generate[T any](fn func(T) []unstructured.Unstructured) Transform[T] {
	return func(items []unstructured.Unstructured, config T) []unstructured.Unstructured {
		return append(items, fn(config)...)
	}
}
