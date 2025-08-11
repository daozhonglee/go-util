package json

import (
	"testing"
)

func TestMarshal(t *testing.T) {
	data := map[string]string{"key": "value"}
	result := MarshalFailSafe(data)
	if result == "" {
		t.Error("Expected non-empty JSON string")
	}
}

func TestClean(t *testing.T) {
	dirty := "Json{\"key\":\"value\"}json"
	clean := Clean(dirty)
	expected := "{\"key\":\"value\"}"
	if clean != expected {
		t.Errorf("Expected %s, got %s", expected, clean)
	}
}

func TestUnmarshal(t *testing.T) {
	jsonStr := `{"key":"value"}`
	var result map[string]string
	err := Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		t.Errorf("Unmarshal failed: %v", err)
	}
	if result["key"] != "value" {
		t.Error("Unmarshal result incorrect")
	}
}
