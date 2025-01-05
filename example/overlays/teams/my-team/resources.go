package myteam

import (
	. "github.com/kimerize/kimerize/example/base/cert-manager"
	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type MyTeamOverlay struct {
	Team        string
	CertManager CertManagerGenerator
}

func (config MyTeamOverlay) Transform(items []unstructured.Unstructured) []unstructured.Unstructured {
	for _, i := range items {
		ls := i.GetLabels()
		if ls == nil {
			ls = map[string]string{}
		}
		ls["team"] = "my-team"
		i.SetLabels(ls)
	}
	return items
}

var DefaultMyTeam = Overlay[MyTeamOverlay]{
	Config: MyTeamOverlay{
		Team:        "my-team",
		CertManager: DefaultCertManager,
	},
}
