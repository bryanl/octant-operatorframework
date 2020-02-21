package component

import (
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func AddStringSectionItem(sections []component.SummarySection, name string, m map[string]interface{}, fields ...string) ([]component.SummarySection, error) {
	text, _, err := unstructured.NestedString(m, fields...)
	if err != nil {
		return nil, err
	}

	return append(sections, component.SummarySection{
		Header:  name,
		Content: component.NewText(text),
	}), nil
}
