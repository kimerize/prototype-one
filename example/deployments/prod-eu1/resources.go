package prodeu1

import (
	. "github.com/kimerize/kimerize/example/base/cert-manager"
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ProdConfig struct {
	Prod string
}

var Resources = ResourceGroup[ProdConfig]{
	Resources: []Generator{
		MyTeam.WithOverrides(func(rg *ResourceGroup[MyTeamConfig]) {
			MustOverride(rg, func(rg *ResourceGroup[CertManagerConfig]) {
				rg.Config.Version = "1.6.2"
			})
		}),
	},
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
	},
}
