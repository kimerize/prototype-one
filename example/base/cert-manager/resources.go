package certmanager

import (
	"fmt"

	. "github.com/kimerize/kimerize/lib"
	. "github.com/kimerize/kimerize/lib/kustomize"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type CertManagerGenerator struct {
	Version string
}

var _ Generator = CertManagerGenerator{}

func (c CertManagerGenerator) Generate() resmap.ResMap {
	return KustomizeBuild(types.Kustomization{
		Resources: []string{
			fmt.Sprintf("https://github.com/cert-manager/cert-manager/releases/download/v%s/cert-manager.yaml", c.Version),
		},
	}, filesys.MakeFsInMemory())
}

var DefaultCertManager = CertManagerGenerator{
	Version: "1.5.0",
}
