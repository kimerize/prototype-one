package myteam

import (
	"fmt"

	. "github.com/kimerize/kimerize/example/base/cert-manager"
	. "github.com/kimerize/kimerize/lib"
)

type MyTeamOverlay struct {
	Team        string
	CertManager CertManagerGenerator
}

func (config MyTeamOverlay) Transform(resources *ResourceList) {
	resources.ErrorAggregator.Append(NewError(fmt.Errorf("MyTeamOverlay.Transform not implemented")))
	resources.ForEach(func(r *Resource) error {
		r.SetLabel("team", config.Team)
		return nil
	})
}

var DefaultMyTeam = Overlay[MyTeamOverlay]{
	Config: MyTeamOverlay{
		Team:        "my-team",
		CertManager: DefaultCertManager,
	},
}
