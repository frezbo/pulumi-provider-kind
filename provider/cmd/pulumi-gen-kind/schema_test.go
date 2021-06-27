package main

import (
	"encoding/json"
	"testing"

	"github.com/alecthomas/jsonschema"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
)

func TestSchema(t *testing.T) {
	schema := jsonschema.Reflect(&v1alpha4.Cluster{})
	bytes, err := json.Marshal(schema)
	if err != nil {
		t.Error(err)
	}
	expected := ""
	if string(bytes) != expected {
		t.Errorf("expected: %s, got: %s", expected, string(bytes))
	}
}
