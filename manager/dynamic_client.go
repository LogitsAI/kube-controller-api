package manager

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// DynamicClient is a wrapper around the client-go dynamic client that looks up
// resources by GroupVersionKind.
type DynamicClient struct {
	dynamic.Interface

	restMapper meta.RESTMapper
}

// NewDynamicClient returns a new DynamicClient.
func NewDynamicClient(mgr manager.Manager) (*DynamicClient, error) {
	dynClient, err := dynamic.NewForConfigAndClient(mgr.GetConfig(), mgr.GetHTTPClient())
	if err != nil {
		return nil, err
	}

	return &DynamicClient{
		Interface:  dynClient,
		restMapper: mgr.GetRESTMapper(),
	}, nil
}

// Kind returns a dynamic client for the resource corresponding to the given GVK.
func (c *DynamicClient) Kind(gvk schema.GroupVersionKind) dynamic.NamespaceableResourceInterface {
	mapping, err := c.restMapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return &errorResource{err: err}
	}
	return c.Resource(mapping.Resource)
}

// errorResource is a dynamic client that always returns an error.
type errorResource struct {
	err error
}

func (r *errorResource) Namespace(namespace string) dynamic.ResourceInterface {
	return r
}

func (r *errorResource) Get(ctx context.Context, name string, options metav1.GetOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, r.err
}

func (r *errorResource) List(ctx context.Context, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	return nil, r.err
}

func (r *errorResource) Create(ctx context.Context, obj *unstructured.Unstructured, options metav1.CreateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, r.err
}

func (r *errorResource) Update(ctx context.Context, obj *unstructured.Unstructured, options metav1.UpdateOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, r.err
}

func (r *errorResource) UpdateStatus(ctx context.Context, obj *unstructured.Unstructured, options metav1.UpdateOptions) (*unstructured.Unstructured, error) {
	return nil, r.err
}

func (r *errorResource) Delete(ctx context.Context, name string, options metav1.DeleteOptions, subresources ...string) error {
	return r.err
}

func (r *errorResource) DeleteCollection(ctx context.Context, options metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return r.err
}

func (r *errorResource) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return nil, r.err
}

func (r *errorResource) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, options metav1.PatchOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, r.err
}

func (r *errorResource) Apply(ctx context.Context, name string, obj *unstructured.Unstructured, opts metav1.ApplyOptions, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, r.err
}

func (r *errorResource) ApplyStatus(ctx context.Context, name string, obj *unstructured.Unstructured, opts metav1.ApplyOptions) (*unstructured.Unstructured, error) {
	return nil, r.err
}
