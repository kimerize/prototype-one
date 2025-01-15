package main

import (
	. "github.com/kimerize/kimerize/example/generator"
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
)

type ProdOverlay struct {
	Prod   string
	MyTeam Overlay[MyTeamOverlay]
}

func (config ProdOverlay) Transform(rl *ResourceList) {
	rl.ForEach(func(r *Resource) error {
		r.SetLabel("prod", config.Prod)
		return nil
	})
}

var MyTeamWithOverrides = DefaultMyTeam.WithOverrides(func(rg *Overlay[MyTeamOverlay]) {
	rg.Config.CertManager.Version = "1.6.2"
})

var Resources = Overlay[MyCorpOverlay]{
	Config: NewMyCorpOverlay(Overlay[ProdOverlay]{
		Config: ProdOverlay{
			Prod:   "prod",
			MyTeam: MyTeamWithOverrides,
		},
	}),

	// prodOverlay.Myeam.CertManager
}
