package wfobj

import (
	"testing"
)

func TestParser(t *testing.T) {
	p, err := NewParser("cube.obj")
	if err != nil {
		t.Fatalf("Unable to create parser. %v", err)
	}
	err = p.Parse()
	if err != nil {
		t.Fatalf("Unable to parse the cube.obj file. %v", err)
	}
}
