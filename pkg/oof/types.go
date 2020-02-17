package oof

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	CatalogSourceGVK = schema.GroupVersionKind{
		Group:   "operators.coreos.com",
		Version: "v1alpha1",
		Kind:    "CatalogSource",
	}
)
