package lib

import (
	"fmt"
	"reflect"
)

type Transformer interface {
	Transform(items *ResourceList)
}

type OverlayTransformer interface {
	Transformer
	SetDefaults()
}

// type Generator interface {
// 	Generate(internal.GenerateOptions) ResourceList
// }

// type Overlay interface {
// 	Generator
// }

type overlay[T any, P interface {
	*T
	OverlayTransformer
}] struct {
	config P
}

// type Overlay[T any] interface {
// 	Generator
// 	Override(override func(*T)) Overlay[T]
// }

func setDefaults(o OverlayTransformer) {
	fmt.Println("setDefaults", o)
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("Overlay config must be a struct")
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		v := field.Addr().Interface()
		if g, ok := v.(OverlayTransformer); ok {
			setDefaults(g)
		}
	}
	o.SetDefaults()
}

func BuildOverlay[T any, P interface {
	*T
	OverlayTransformer
}](override ...func(*T)) ResourceList {
	t := *new(T)
	var p P = &t
	setDefaults(p)
	for _, o := range override {
		o(&t)
	}
	overlay := overlay[T, P]{
		config: p,
	}
	return overlay.Generate()
}

type dummyOverlayConfig struct {
	Kind string
}

func (dummyOverlayConfig) Transform(items *ResourceList) {}

func (d *dummyOverlayConfig) SetDefaults() {
	d.Kind = "dummy"
}

var _ OverlayTransformer = &dummyOverlayConfig{}

// var _ Generator = BuildOverlay[dummyOverlayConfig]()

// func (t *overlay[T, P]) Override(f func(*T)) Overlay[T] {
// 	f(t.config)
// 	return t
// }

// func (t *overlay[T, P]) Generate() ResourceList {
// 	v := reflect.ValueOf(t.config)
// 	if v.Kind() == reflect.Ptr {
// 		v = v.Elem()
// 	}
// 	if v.Kind() != reflect.Struct {
// 		panic("Overlay config must be a struct")
// 	}

// 	result := NewResourceList()
// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)
// 		if g, ok := field.Interface().(OverlayConfig); ok {
// 			// if d, ok := reflect.ValueOf(&g).Interface().(Defaulter); ok {
// 			// 	d.SetDefaults()
// 			// }
// 			result.Absorb(g.Generate())
// 		}
// 	}
// 	t.config.Transform(result)
// 	return *result
// }

func (t overlay[T, P]) Generate() ResourceList {
	fmt.Println("generate", t.config)
	return generate(t.config)
}

func generate(o OverlayTransformer) ResourceList {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic("Overlay config must be a struct")
	}

	result := NewResourceList()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		v := field.Addr().Interface()
		if g, ok := v.(OverlayTransformer); ok {
			result.Absorb(generate(g))
		}
	}
	o.Transform(result)
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
