package wfobj

import (
	"testing"
)

func TestMeshLoader(t *testing.T) {
	for _, test := range testdata {
		if test.ignore {
			continue
		}
		t.Logf("Title: %v", test.title)
		p := NewLiteralParser(test.objlit)
		p.Debug = &PrintState{}
		go p.Parse()

		m, err := LoadMesh(p.Tokens)
		if err != nil {
			t.Fatalf("Unable to load mesh: %v", err)
		}
		if len(m.Faces) != len(test.mesh.Faces) {
			t.Fatalf("Mesh should have 2 faces")
			continue
		}
		for i, _ := range m.Faces {
			if !m.Faces[i].Same(&test.mesh.Faces[i]) {
				t.Fatalf("Faces are different. Expecting %v got %v", test.mesh.Faces[i], m.Faces[i])
			}
			if !m.Faces[i].Normals.Same(test.mesh.Faces[i].Normals) {
				t.Fatalf("Faces normals are different. Expecting %v got %v", test.mesh.Faces[i].Normals,
					m.Faces[i].Normals)
			}
		}
	}
}
