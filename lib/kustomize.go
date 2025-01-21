package lib

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func init() {
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
}

func KustomizeBuild(kustomize types.Kustomization, fs filesys.FileSystem) ResourceList {
	options := krusty.MakeDefaultOptions()
	options.PluginConfig.HelmConfig.Enabled = true
	options.PluginConfig.HelmConfig.Command = "helm"
	k := krusty.MakeKustomizer(options)

	kBytes, err := yaml.Marshal(kustomize)
	if err != nil {
		// TODO:
	}
	err = fs.WriteFile("kustomization.yaml", kBytes)
	if err != nil {
		// TODO:
	}

	rm, err := k.Run(fs, ".")
	if err != nil {
		// TODO:
		fmt.Println("kustomize Error:", err)
		// return false, err
	}
	result := ResourceList{}
	for _, r := range rm.Resources() {
		result.Append(ResourceFrom(r.RNode))
	}
	return result
}
