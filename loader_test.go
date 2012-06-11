package wfobj

import (
	"testing"
)

func TestMeshLoader(t *testing.T) {
	p := NewLiteralParser(objlit)
	p.Debug = &PrintState{}
	go p.Parse()

	m, err := LoadMesh(p.Tokens)
	if err != nil {
		t.Fatalf("Unable to load mesh: %v", err)
	}
	if len(m.Faces) != 2 {
		t.Fatalf("Mesh should have 2 faces")
	}
	for i, _ := range m.Faces {
		if !m.Faces[i].Same(&model.Faces[i]) {
			t.Fatalf("Faces are different. Expecting %v got %v", model.Faces[i], m.Faces[i])
		}
	}
}
