package lib

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Transform[T any] func(items []unstructured.Unstructured, config T) []unstructured.Unstructured

type Generator interface {
	Generate() []unstructured.Unstructured
}

type ResourceGroup[T any] struct {
	Resources []Generator
	Transform[T]
	Config T
}

var _ Generator = ResourceGroup[any]{}

func (t ResourceGroup[T]) Generate() []unstructured.Unstructured {
	items := []unstructured.Unstructured{}
	for _, r := range t.Resources {
		items = r.Generate()
	}
	items = t.Transform(items, t.Config)
	return items
}

func (t ResourceGroup[T]) WithOverrides(override func(*ResourceGroup[T])) ResourceGroup[T] {
	override(&t)
	return t
}

func MustOverride[T any, K any](rg *ResourceGroup[T], override func(*ResourceGroup[K])) {
	found := false
	for i, r := range rg.Resources {
		if t, ok := r.(ResourceGroup[K]); ok {
			found = true
			override(&t)
			rg.Resources[i] = t
		}
	}
	if !found {
		panic(fmt.Errorf("transformer not found"))
	}
}

func Generate[T any](fn func(T) []unstructured.Unstructured) Transform[T] {
	return func(items []unstructured.Unstructured, config T) []unstructured.Unstructured {
		return append(items, fn(config)...)
	}
}
