package controllerpb

import "k8s.io/apimachinery/pkg/runtime/schema"

func (gvk *GroupVersionKind) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   gvk.Group,
		Version: gvk.Version,
		Kind:    gvk.Kind,
	}
}
