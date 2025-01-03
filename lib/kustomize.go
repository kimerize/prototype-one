package lib

import (

	// "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func KustomizeBuild(kustomize types.Kustomization, fs filesys.FileSystem) []unstructured.Unstructured {
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

	m, err := k.Run(fs, ".")
	if err != nil {
		// TODO:
		// return false, err
	}
	result := []unstructured.Unstructured{}
	for _, r := range m.Resources() {
		// TODO: handle error
		m, _ := r.Map()
		result = append(result, unstructured.Unstructured{Object: m})
	}
	return result
}
