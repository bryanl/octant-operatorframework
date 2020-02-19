package oof_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func requireJSONEq(t *testing.T, expected interface{}, got interface{}) {
	e, err := json.Marshal(expected)
	require.NoError(t, err)

	w, err := json.Marshal(got)
	require.NoError(t, err)

	require.JSONEq(t, string(e), string(w))
}

func loadObject(t *testing.T, name string) *unstructured.Unstructured {
	f, err := os.Open(filepath.Join("testdata", name))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, f.Close())
	}()

	var m map[string]interface{}
	require.NoError(t, json.NewDecoder(f).Decode(&m))

	return &unstructured.Unstructured{Object: m}
}
