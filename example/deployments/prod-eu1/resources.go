package main

import (
	. "github.com/kimerize/kimerize/example/generator"
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
	"sigs.k8s.io/kustomize/api/resmap"
)

type ProdOverlay struct {
	Prod   string
	MyTeam Overlay[MyTeamOverlay]
}

func (config ProdOverlay) Transform(rm resmap.ResMap) {
	for _, r := range rm.Resources() {
		labels := r.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels["prod"] = config.Prod
		r.SetLabels(labels)
	}
}

var Resources = Overlay[MyCorpOverlay]{
	Config: NewMyCorpOverlay(Overlay[ProdOverlay]{
		Config: ProdOverlay{
			Prod: "prod",
			MyTeam: DefaultMyTeam.WithOverrides(func(rg *Overlay[MyTeamOverlay]) {
				rg.Config.CertManager.Version = "1.6.2"
			}),
		},
	}),
}
