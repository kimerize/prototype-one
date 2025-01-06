package lib

import (
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func KustomizeBuild(kustomize types.Kustomization, fs filesys.FileSystem) resmap.ResMap {
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
	return m
}
