package main

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/k8s-digester/pkg/resolve"
	myteam "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type ProdOverlay struct {
	Prod   string
	MyTeam myteam.MyTeamOverlay
}

var _ Overlay = &ProdOverlay{}

func (config *ProdOverlay) Transform(rl *ResourceList) {
	rl.ForEach(func(r *Resource) {
		r.SetLabel("prod", config.Prod)
	})
}

func (config *ProdOverlay) SetDefaults() {
	config.Prod = "prod"
}

var Resources ResourceList = func() ResourceList {
	rl := BuildOverlay(func(po *ProdOverlay) {
		po.MyTeam.CertManager.Version = "1.5.0"
	})
	rl.ForEach(func(r *Resource) {
		ModifyAs(r, func(r *yaml.RNode) {
			resolve.ImageTags(context.TODO(), logr.Discard(), nil, r, nil)
		})
	})
	return rl
}()

var Publisher PackagePublisher = KustomizePublisher()
