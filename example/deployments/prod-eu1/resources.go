package main

import (
	. "github.com/kimerize/kimerize/example/base/cert-manager"
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

func (config *ProdOverlay) SetDefaults() {
	config.Prod = "prod"
}

var Resources = NewOverlay[ProdOverlay]().
	Override(func(po *ProdOverlay) {
		po.MyTeam.Override(func(mto *MyTeamOverlay) {
			mto.CertManager.Override(func(cm *CertManager) {
				cm.Version = "1.5.0"
			})
		})
	})
