package myteam

import (
	. "github.com/kimerize/kimerize/example/base/cert-manager"
	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type MyTeamConfig struct {
	Team string
}

var MyTeam = ResourceGroup[MyTeamConfig]{
	Resources: []Generator{CertManager},
	Transform: func(items []unstructured.Unstructured, config MyTeamConfig) []unstructured.Unstructured {
		for _, i := range items {
			ls := i.GetLabels()
			if ls == nil {
				ls = map[string]string{}
			}
			ls["team"] = "my-team"
			i.SetLabels(ls)
		}
		return items
	},
}
