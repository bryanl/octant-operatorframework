package oof

import "k8s.io/apimachinery/pkg/runtime/schema"

var (
	CatalogSourceGVK = schema.GroupVersionKind{
		Group:   "operators.coreos.com",
		Version: "v1alpha1",
		Kind:    "CatalogSource",
	}
	SubscriptionGVK = schema.GroupVersionKind{
		Group:   "operators.coreos.com",
		Version: "v1alpha1",
		Kind:    "Subscription",
	}
)

const (
	ActionSubscribeToPackage = "oof.bryanl.dev/subscribe-to-package"
)
