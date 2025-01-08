package lib

import (
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/kustomize/v5/commands/build"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

type KustomizeWriter struct {
}

func Write(dir string, items resmap.ResMap) {
	writer := build.MakeWriter(filesys.FileSystemOrOnDisk{})
	writer.WriteIndividualFiles(dir, items)
}
