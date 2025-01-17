package lib

import (
	"reflect"
)

type Defaulter interface {
	SetDefaults()
}

type OverlayInterface[T Transformer] interface {
	Generator
}

type Transformer interface {
	Transform(items *ResourceList)
}

type OverlayConfig interface {
	Transformer
	Defaulter
}

type Generator interface {
	Generate() ResourceList
}

type overlay[T any, P interface {
	*T
	OverlayConfig
}] struct {
	config P
}

type Overlay[T any] interface {
	Generator
	Override(override func(*T)) Overlay[T]
}

func NewOverlay[T any, P interface {
	*T
	OverlayConfig
}]() Overlay[T] {
	t := *new(T)
	var p P = &t
	P.SetDefaults(&t)
	overlay := overlay[T, P]{
		config: p,
	}
	return &overlay
}

type dummyOverlayConfig struct {
	kind string
}

func (dummyOverlayConfig) Transform(items *ResourceList) {}

func (d *dummyOverlayConfig) SetDefaults() {
	d.kind = "dummy"
}

var _ OverlayConfig = &dummyOverlayConfig{}

var _ Generator = NewOverlay[dummyOverlayConfig]()

func (t *overlay[T, P]) Override(f func(*T)) Overlay[T] {
	f(t.config)
	return t
}

func (t *overlay[T, P]) Generate() ResourceList {
	v := reflect.ValueOf(t.config)
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
			if d, ok := reflect.ValueOf(&g).Interface().(Defaulter); ok {
				d.SetDefaults()
			}
			result.Absorb(g.Generate())
		}
	}
	t.config.Transform(result)
	return *result
}

// func (t Overlay[T]) WithOverrides(override func(*Overlay[T])) Overlay[T] {
// 	override(&t)
// 	return t
// }

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
