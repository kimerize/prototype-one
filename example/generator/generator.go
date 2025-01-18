package generator

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/k8s-digester/pkg/resolve"
	. "github.com/kimerize/kimerize/lib"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type MyCorpOverlay struct {
	Generator
}

var _ Transformer = MyCorpOverlay{}

func NewMyCorpOverlay(g Generator) MyCorpOverlay {
	return MyCorpOverlay{Generator: g}
}

// Generate implements lib.Generator.
func (m MyCorpOverlay) Transform(rl *ResourceList) {
	rl.ForEach(func(r *Resource) {
		ModifyAs(r, func(r *yaml.RNode) {
			resolve.ImageTags(context.TODO(), logr.Discard(), nil, r, nil)
		})
	})
}

var _ Generator = MyCorpOverlay{}
