package certmanager

import (
	"fmt"

	. "github.com/kimerize/kimerize/lib"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type CertManagerConfig struct {
	Version string
}

var CertManager = ResourceGroup[CertManagerConfig]{
	Transform: Generate(func(c CertManagerConfig) []unstructured.Unstructured {
		return KustomizeBuild(types.Kustomization{
			Resources: []string{
				fmt.Sprintf("https://github.com/cert-manager/cert-manager/releases/download/v%s/cert-manager.yaml", c.Version),
			},
		}, filesys.MakeFsInMemory())
	}),
	Config: CertManagerConfig{
		Version: "1.5.0",
	},
}
