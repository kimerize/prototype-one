package main

import (
	. "github.com/kimerize/kimerize/example/overlays/teams/my-team"
	. "github.com/kimerize/kimerize/lib"
)

type ProdOverlay struct {
	Prod   string
	MyTeam MyTeamOverlay
}

func (config *ProdOverlay) Transform(rl *ResourceList) {
	rl.ForEach(func(r *Resource) {
		r.SetLabel("prod", config.Prod)
	})
}

func (config *ProdOverlay) SetDefaults() {
	config.Prod = "prod"
}

var Resources ResourceList = func() ResourceList {
	rl := BuildOverlay(func(po *ProdOverlay) {
		po.MyTeam.CertManager.Version = "1.5.0"
	})
	rl.ForEach(func(r *Resource) {
		r.SetLabel("cluster", "my-cluster")
	})
	return rl
}()

var Publisher PackagePublisher = KustomizePublisher()
