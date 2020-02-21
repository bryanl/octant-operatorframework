package oof

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func toUnstructured(in runtime.Object) (*unstructured.Unstructured, error) {
	if in == nil {
		return nil, fmt.Errorf("object is nil")
	}

	m, err := runtime.DefaultUnstructuredConverter.ToUnstructured(in)
	if err != nil {
		return nil, err
	}

	return &unstructured.Unstructured{Object: m}, nil
}
