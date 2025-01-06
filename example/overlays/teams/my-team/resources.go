package myteam

import (
	. "github.com/kimerize/kimerize/example/base/cert-manager"
	. "github.com/kimerize/kimerize/lib"
	"sigs.k8s.io/kustomize/api/resmap"
)

type MyTeamOverlay struct {
	Team        string
	CertManager CertManagerGenerator
}

func (config MyTeamOverlay) Transform(rm resmap.ResMap) {
	for _, r := range rm.Resources() {
		labels := r.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}
		labels["team"] = "my-team"
		r.SetLabels(labels)
	}
}

var DefaultMyTeam = Overlay[MyTeamOverlay]{
	Config: MyTeamOverlay{
		Team:        "my-team",
		CertManager: DefaultCertManager,
	},
}
