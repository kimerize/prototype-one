package main

import (
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ProdTransformer struct {
	Prod   string
	MyTeam ResourceGroup[MyTeamTransformer]
}

func (config ProdTransformer) Transform(items []unstructured.Unstructured) []unstructured.Unstructured {
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

var Resources = ResourceGroup[ProdTransformer]{
	Config: ProdTransformer{
		Prod: "prod",
		MyTeam: MyTeam.WithOverrides(func(rg *ResourceGroup[MyTeamTransformer]) {
			rg.Config.CertManager.Version = "1.6.2"
		}),
	},
}
