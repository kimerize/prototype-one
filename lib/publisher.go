package lib

import (
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/kustomize/v5/commands/build"
	"sigs.k8s.io/kustomize/kyaml/filesys"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Package struct {
	Resources ResourceList
	*packages.Package
}

func (p Package) localPathOutput() string {
	pkgRelPath, err := filepath.Rel(p.Module.Dir, p.Dir)
	if err != nil {
		panic(err)
	}
	return filepath.Join(p.Module.Dir, "zz_generated", pkgRelPath)
}

type PackagePublisher interface {
	Publish(Package) error
}

type kustomizePublisher struct {
}

func KustomizePublisher() PackagePublisher {
	return kustomizePublisher{}
}

var _ PackagePublisher = kustomizePublisher{}

// Publish implements PackagePublisher.
func (k kustomizePublisher) Publish(p Package) error {
	rm := resmap.New()
	p.Resources.ForEach(func(r *Resource) {
		rm.Append(&resource.Resource{
			RNode: *r.rnode.Copy(),
		})
	})
	outputDir := p.localPathOutput()
	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return err
	}
	writer := build.MakeWriter(filesys.FileSystemOrOnDisk{})
	return writer.WriteIndividualFiles(outputDir, rm)
}

type localPackagePublisher struct {
}

// Publish implements PackagePublisher.
func (l localPackagePublisher) Publish(p Package) error {
	writer := kio.LocalPackageWriter{
		PackagePath: p.localPathOutput(),
	}
	nodes := []*yaml.RNode{}
	p.Resources.ForEach(func(r *Resource) {
		nodes = append(nodes, r.rnode.Copy())
	})
	return writer.Write(nodes)
}

func LocalPackageWriter() localPackagePublisher {
	return localPackagePublisher{}
}

var _ PackagePublisher = localPackagePublisher{}
