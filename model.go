package wfobj

// Represent a 3D vertex
type Vertex struct {
	X, Y, Z float32
}

// Check if two vertexes are the same
func (v *Vertex) Same(other *Vertex) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// Represent a vertex list
type VertexList []Vertex

// Represent one face of the object
// Vertices must be in the right draw order
type Face struct {
	Vertices VertexList
	Normals  VertexList
}

// Check if two faces are equal
// ie, same vertices in the same order
func (f *Face) Same(other *Face) bool {
	if len(f.Vertices) != len(other.Vertices) {
		return false
	}
	for i, _ := range f.Vertices {
		if !f.Vertices[i].Same(&other.Vertices[i]) {
			return false
		}
	}
	return true
}

// Represent a mesh made by a collection of faces/material
// TODO Implement material
type Mesh struct {
	Faces []Face
}
