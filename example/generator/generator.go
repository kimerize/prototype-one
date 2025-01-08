package generator

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/k8s-digester/pkg/resolve"
	. "github.com/kimerize/kimerize/lib"
	"sigs.k8s.io/kustomize/api/resmap"
)

type MyCorpOverlay struct {
	Generator
}

var _ Transformer = MyCorpOverlay{}

func NewMyCorpOverlay(g Generator) MyCorpOverlay {
	return MyCorpOverlay{Generator: g}
}

// Generate implements lib.Generator.
func (m MyCorpOverlay) Transform(rm resmap.ResMap) {
	for _, r := range rm.Resources() {
		resolve.ImageTags(context.TODO(), logr.Discard(), nil, &r.RNode, nil)
	}
}

var _ Generator = MyCorpOverlay{}
