package lib

import "sigs.k8s.io/kustomize/kyaml/yaml"

type Resource struct {
	yaml.RNode
}

type ResourceList interface {
	Append(Resource)
	AppendAll(ResourceList)
	Resources() []Resource
}

type resourceListImpl struct {
	resources []Resource
}

// Resources implements ResourceList.
func (r *resourceListImpl) Resources() []Resource {
	panic("unimplemented")
}

// Append implements ResourceList.
func (r *resourceListImpl) Append(Resource) {
	panic("unimplemented")
}

// AppendAll implements ResourceList.
func (r *resourceListImpl) AppendAll(ResourceList) {
	panic("unimplemented")
}

func NewResourceList() ResourceList {
	return &resourceListImpl{}
}
