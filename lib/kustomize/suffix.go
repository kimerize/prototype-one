package kustomize

import (
	"encoding/json"
	"reflect"
	"regexp"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/kustomize/api/hasher"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kyaml/openapi"
	"sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/kustomize/kyaml/yaml/walk"
)

var h = &hasher.Hasher{}

// func ModifyHashSuffixedResource(rm resmap.ResMap, name string, fn func(*yaml.RNode)) {
// 	for _, r := range rm.Resources() {
// 		copy := r.Copy()
// 		re := regexp.MustCompile(`(^.+)-([a-z0-9]+)$`)
// 		if match := re.FindStringSubmatch(copy.GetName()); match != nil {
// 			if name != match[1] {
// 				continue
// 			}
// 			copy.SetName(name)
// 			wantedHash := match[2]
// 			if gotHash, _ := h.Hash(copy); wantedHash == gotHash {
// 				fn(&r.RNode)
// 				newHash, _ := h.Hash(&r.RNode)
// 				r.SetName(name + "-" + newHash)
// 			}
// 			return
// 		}
// 	}
// }

func ModifyAs[T any](r *yaml.RNode, fn func(*T)) {
	var t T

	b, err := r.MarshalJSON()
	if err != nil {
		// TODO:
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		// TODO:
	}

	fn(&t)

	b, err = json.Marshal(t)
	if err != nil {
		// TODO:
	}
	modified, err := yaml.Parse(string(b))
	if err != nil {
		// TODO:
	}

	r.SetYNode(modified.YNode())
}

type valueSetterWalker struct {
	value    string
	newValue string
}

func (n *valueSetterWalker) visitCollection(elements []*yaml.RNode) error {
	for _, element := range elements {
		_, err := walk.Walker{
			Visitor: n,
			Sources: walk.Sources{element},
		}.Walk()
		if err != nil {
			return err
		}
	}
	return nil
}

// VisitList implements walk.Visitor.
func (n *valueSetterWalker) VisitList(sources walk.Sources, _ *openapi.ResourceSchema, _ walk.ListKind) (*yaml.RNode, error) {
	elements, error := sources.Dest().Elements()
	if error != nil {
		return nil, error
	}
	if err := n.visitCollection(elements); err != nil {
		return nil, err
	}
	return sources.Dest(), nil
}

// VisitMap implements walk.Visitor.
func (n *valueSetterWalker) VisitMap(sources walk.Sources, _ *openapi.ResourceSchema) (*yaml.RNode, error) {
	fields, error := sources.Dest().FieldRNodes()
	if error != nil {
		return nil, error
	}
	if err := n.visitCollection(fields); err != nil {
		return nil, err
	}
	return sources.Dest(), nil
}

// VisitScalar implements walk.Visitor.
func (n *valueSetterWalker) VisitScalar(sources walk.Sources, _ *openapi.ResourceSchema) (*yaml.RNode, error) {
	return yaml.ValueReplacer{StringMatch: n.value, Replace: n.newValue}.Filter(sources.Dest())
}

var _ walk.Visitor = &valueSetterWalker{}

// TODO: add namespace
func ModifyHashSuffixedResource[T any](rm resmap.ResMap, name types.NamespacedName, fn func(*T)) {
	kind := reflect.TypeOf((*T)(nil)).Elem().Name()
	var foundName string
	var newName string
	for _, r := range rm.Resources() {
		if r.GetKind() != kind {
			continue
		}
		copy := r.Copy()
		re := regexp.MustCompile(`(^.+)-([a-z0-9]{10})$`)
		if match := re.FindStringSubmatch(copy.GetName()); match != nil {
			if name.Name != match[1] || r.GetNamespace() != name.Namespace {
				continue
			}
			copy.SetName(name.Name)
			wantedHash := match[2]
			if gotHash, _ := h.Hash(copy); wantedHash == gotHash {
				foundName = r.GetName()
				ModifyAs(&r.RNode, fn)
				newHash, _ := h.Hash(&r.RNode)
				r.SetName(name.Name + "-" + newHash)
				newName = r.GetName()
			}
			break
		}
	}

	for _, r := range rm.Resources() {
		if r.GetKind() == kind && r.GetName() == newName {
			continue
		}
		_, err := walk.Walker{
			Visitor: &valueSetterWalker{value: foundName, newValue: newName},
			Sources: walk.Sources{&r.RNode},
		}.Walk()
		if err != nil {
			// TODO: handle error
		}

	}
}
