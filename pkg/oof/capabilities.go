package oof

import (
	"github.com/vmware-tanzu/octant/pkg/plugin"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func InitCapabilities() *plugin.Capabilities {
	c := plugin.Capabilities{
		SupportsPrinterConfig: []schema.GroupVersionKind{
			CatalogSourceGVK,
		},
		IsModule: true,
	}

	return &c
}
