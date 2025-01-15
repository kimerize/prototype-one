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
	addError(resources)
	resources.ForEach(func(r *Resource) error {
		r.SetLabel("team", config.Team)
		return nil
	})
}

func addError(resources *ResourceList) {
	fmt.Errorf("MyTeamOverlay.Transform not implemented")
	resources.Error(fmt.Errorf("MyTeamOverlay.Transform not implemented"))
}

var DefaultMyTeam = Overlay[MyTeamOverlay]{
	Config: MyTeamOverlay{
		Team:        "my-team",
		CertManager: DefaultCertManager,
	},
}
