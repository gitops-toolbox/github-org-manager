package githuborg

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	result, err := LoadConfig("../repos/test")

	if err != nil {
		t.Fatal(err)
	}

	templator, ok := result["templator"]

	if !ok {
		t.Fatal("Result does not contian expected templator repo")
	}

	if *templator.Name != "templator" {
		t.Fatalf("Expected templator repo, got %s", *templator.Name)
	}
}
