package wfobj

// This is a invalid 3d object
// but it's a syntatic  valid .obj file
// so for the purpose of these testing
// this just works
var objlit = `# comment
v 1.0 1.0 1.0
v 0.0 1.0 0.0
f 1 2
f 2 1
`

var model = &Mesh{
	[]Face{
		Face{VertexList{Vertex{1.0, 1.0, 1.0}, Vertex{0.0, 1.0, 0.0}}},
		Face{VertexList{Vertex{0.0, 1.0, 0.0}, Vertex{1.0, 1.0, 1.0}}},
	},
}

type PrintState struct{}

func (ps *PrintState) State(p *Parser) {
	print("STATE: ", p.String(), "\n")
}

func (ps *PrintState) Emit(t *Token) {
	print("EMIT: ", t.String(), "\n")
}
