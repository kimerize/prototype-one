package main

import (
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PordOverlay struct {
	Prod   string
	MyTeam Overlay[MyTeamOverlay]
}

func (config PordOverlay) Transform(items []unstructured.Unstructured) []unstructured.Unstructured {
	for _, i := range items {
		ls := i.GetLabels()
		if ls == nil {
			ls = map[string]string{}
		}
		ls["prod"] = config.Prod
		i.SetLabels(ls)
	}
	return items
}

var Resources = Overlay[PordOverlay]{
	Config: PordOverlay{
		Prod: "prod",
		MyTeam: MyTeam.WithOverrides(func(rg *Overlay[MyTeamOverlay]) {
			rg.Config.CertManager.Version = "1.6.2"
		}),
	},
}
