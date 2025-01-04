package main

import (
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ProdConfig struct {
	Prod   string
	MyTeam ResourceGroup[MyTeamConfig]
}

var Resources = ResourceGroup[ProdConfig]{
	Transform: func(items []unstructured.Unstructured, config ProdConfig) []unstructured.Unstructured {
		for _, i := range items {
			ls := i.GetLabels()
			if ls == nil {
				ls = map[string]string{}
			}
			ls["prod"] = config.Prod
			i.SetLabels(ls)
		}
		return items
	},
	Config: ProdConfig{
		Prod: "prod",
		MyTeam: MyTeam.WithOverrides(func(rg *ResourceGroup[MyTeamConfig]) {
			rg.Config.CertManager.Config.Version = "1.6.2"
		}),
	},
}
