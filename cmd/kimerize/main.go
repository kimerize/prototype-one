package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	"github.com/go-git/go-git/v5"
	"github.com/kimerize/kimerize/lib"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func main() {
	target, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %v\n", err)
		return
	}
	if len(os.Args) > 1 {
		if filepath.IsAbs(os.Args[1]) {
			target = os.Args[1]
		} else {
			target = filepath.Join(target, os.Args[1])
		}
	}
	fmt.Println("Processing directory:", target)

	// Find the git repository root starting from the target directory
	r, err := git.PlainOpenWithOptions(target, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		fmt.Printf("Error opening git repository: %v\n", err)
		return
	}
	var goDirs []string
	err = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			matches, err := filepath.Glob(filepath.Join(path, "*.go"))
			if err != nil {
				return err
			}
			if len(matches) > 0 {
				goDirs = append(goDirs, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory tree: %v\n", err)
		return
	}
	packages, err := loader.LoadRoots(goDirs...)
	if err != nil {
		fmt.Printf("Error loading roots: %v\n", err)
		return
	}

	// TODO: find a module root and use that as the root directory
	wt, err := r.Worktree()
	if err != nil {
		fmt.Printf("Error getting worktree: %v\n", err)
		return
	}
	rootDir := wt.Filesystem.Root()
	processPackages(packages, rootDir)
}

func processPackages(packages []*loader.Package, rootDir string) error {
	for _, p := range packages {
		if p.Package.Name != "main" {
			continue
		}
		fmt.Printf("Package: %s --- %s\n", p.PkgPath, p.Dir)
		if err := processPackage(p, rootDir); err != nil {
			fmt.Printf("Error processing package %s: %v\n", p.PkgPath, err)
			continue
		}
	}
	return nil
}

func buildPlugin(pkg *loader.Package) (*plugin.Plugin, error) {
	pluginPath := filepath.Join(pkg.Dir, "plugin.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginPath, pkg.Dir)
	cmd.Dir = pkg.Dir

	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("build error: %v\n%s", err, out)
	}

	return plugin.Open(pluginPath)
}

func getGenerator(pl *plugin.Plugin) (lib.Generator, error) {
	symbol, err := pl.Lookup("Resources")
	if err != nil {
		return nil, err
	}

	generator, ok := symbol.(lib.Generator)
	if !ok {
		return nil, fmt.Errorf("unexpected function signature")
	}
	return generator, nil
}

func generateNodes(generator lib.Generator) ([]*yaml.RNode, error) {
	var nodes []*yaml.RNode
	for _, r := range generator.Generate() {
		n, err := yaml.FromMap(r.Object)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}

func processPackage(pkg *loader.Package, rootDir string) error {
	tmpDir, err := os.MkdirTemp("", "plugin-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	pluginPath := filepath.Join(tmpDir, "plugin.so")
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", pluginPath, pkg.Dir)
	cmd.Dir = pkg.Dir

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build error: %v\n%s", err, out)
	}

	pl, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}

	generator, err := getGenerator(pl)
	if err != nil {
		return err
	}

	nodes, err := generateNodes(generator)
	if err != nil {
		return err
	}

	relPath, _ := filepath.Rel(rootDir, pkg.Dir)
	outputPath := filepath.Join(rootDir, "zz_generated", relPath)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	writer := kio.LocalPackageWriter{
		PackagePath: outputPath,
	}
	return writer.Write(nodes)
}
