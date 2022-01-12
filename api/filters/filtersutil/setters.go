package filtersutil

import (
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// SetFn is a function that accepts an RNode to possibly modify.
type SetFn func(*yaml.RNode) error

// SetScalar returns a SetFn to set a scalar value
func SetScalar(value string) SetFn {
	return func(node *yaml.RNode) error {
		return node.PipeE(yaml.FieldSetter{StringValue: value})
	}
}

// SetEntry returns a SetFn to set an entry in a map
func SetEntry(key, value, tag string) SetFn {
	n := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: value,
		Tag:   tag,
	}
	if tag == yaml.NodeTagString && yaml.IsYaml1_1NonString(n) {
		n.Style = yaml.DoubleQuotedStyle
	}
	return func(node *yaml.RNode) error {
		return node.PipeE(yaml.FieldSetter{
			Name:  key,
			Value: yaml.NewRNode(n),
		})
	}
}

type Setter struct {
	// SetScalarCallback will be invoked for each call to SetScalar
	SetScalarCallback func(value string, node *yaml.RNode)
	// SetEntryCallback will be invoked for each call to SetEntry
	SetEntryCallback func(key, value, tag string, node *yaml.RNode)
}

func (s Setter) SetScalar(value string) SetFn {
	origSetScalar := SetScalar(value)
	return func(node *yaml.RNode) error {
		if s.SetScalarCallback != nil {
			s.SetScalarCallback(value, node)
		}
		return origSetScalar(node)
	}
}

func (s Setter) SetEntry(key, value, tag string) SetFn {
	origSetEntry := SetEntry(key, value, tag)
	return func(node *yaml.RNode) error {
		if s.SetEntryCallback != nil {
			s.SetEntryCallback(key, value, tag, node)
		}
		return origSetEntry(node)
	}
}
