package wfobj

// Represent a 3D vertex
type Vertex struct {
	X, Y, Z float32
}

// Check if two vertexes are the same
func (v *Vertex) Same(other *Vertex) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}

// Subtract one for another
func (v *Vertex) Sub(other *Vertex) (ret *Vertex) {
	ret = &Vertex{other.X - v.X, other.Y - v.Y, other.Z - v.Z}
	return
}

// Add on vector to another
func (v *Vertex) Add(other *Vertex) (ret *Vertex) {
	ret = &Vertex{other.X + v.X, other.Y + v.Y, other.Z + v.Z}
	return
}

// Represent a vertex list
type VertexList []Vertex

func (v VertexList) Same(other VertexList) bool {
	if len(v) != len(other) {
		return false
	}
	for i, _ := range other {
		if !v[i].Same(&other[i]) {
			return false
		}
	}
	return true
}

// Represent one face of the object
// Vertices must be in the right draw order
type Face struct {
	Vertices VertexList
	Normals  VertexList
}

// Check if two faces are equal
// ie, same vertices in the same order
func (f *Face) Same(other *Face) bool {
	return f.Vertices.Same(other.Vertices)
}

// Represent a mesh made by a collection of faces/material
// TODO Implement material
type Mesh struct {
	Faces []Face
}
