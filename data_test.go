package wfobj

// This is a invalid 3d object
// but it's a syntatic  valid .obj file
// so for the purpose of these testing
// this just works

type TestData struct {
	title  string
	ignore bool
	mesh   *Mesh
	tokens []Token
	objlit string
}

var testdata = []TestData{
	{
		title:  "Simple mesh",
		ignore: false,
		mesh: &Mesh{
			[]Face{
				Face{VertexList{Vertex{1.0, 1.0, 1.0}, Vertex{0.0, 1.0, 0.0}}, VertexList{}},
				Face{VertexList{Vertex{0.0, 1.0, 0.0}, Vertex{1.0, 1.0, 1.0}}, VertexList{}},
			},
		},
		objlit: `# comment
v 1.0 1.0 1.0
v 0.0 1.0 0.0
f 1 2
f 2 1
`,
		tokens: []Token{

			// Vertex
			Token{"", VertexDecl, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},

			// Vertex
			Token{"", VertexDecl, Position{}},
			Token{"0.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"0.0", NumberLit, Position{}},

			// Face	
			Token{"", FaceDecl, Position{}},
			Token{"1", NumberLit, Position{}},
			Token{"2", NumberLit, Position{}},

			// Face	
			Token{"", FaceDecl, Position{}},
			Token{"2", NumberLit, Position{}},
			Token{"1", NumberLit, Position{}},

			Token{"", Eof, Position{}},
		},
	},

	// Mesh with normals
	{
		title:  "Mesh with normals",
		ignore: false,
		mesh: &Mesh{
			[]Face{
				Face{VertexList{Vertex{1.0, 1.0, 1.0}, Vertex{0.0, 1.0, 0.0}}, VertexList{
					Vertex{1.0, 1.0, 1.0}, Vertex{0.0, 1.0, 0.0},
				}},
				Face{VertexList{Vertex{0.0, 1.0, 0.0}, Vertex{1.0, 1.0, 1.0}}, VertexList{
					Vertex{0.0, 1.0, 0.0}, Vertex{1.0, 1.0, 1.0},
				}},
			},
		},
		objlit: `# comment
v 1.0 1.0 1.0
v 0.0 1.0 0.0
vn 1.0 1.0 1.0
vn 0.0 1.0 0.0
f 1//1 2//2
f 2//2 1//1
`,
		tokens: []Token{

			// Vertex
			Token{"", VertexDecl, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},

			// Vertex
			Token{"", VertexDecl, Position{}},
			Token{"0.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"0.0", NumberLit, Position{}},

			// Normal
			Token{"", NormalDecl, Position{}},
			Token{"0.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"0.0", NumberLit, Position{}},

			// Normal
			Token{"", NormalDecl, Position{}},
			Token{"0.0", NumberLit, Position{}},
			Token{"1.0", NumberLit, Position{}},
			Token{"0.0", NumberLit, Position{}},

			// Face	
			Token{"", FaceDecl, Position{}},
			Token{"1", NumberLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"1", NumberLit, Position{}},

			// vector
			Token{"2", NumberLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"2", NumberLit, Position{}},

			// Face	
			Token{"", FaceDecl, Position{}},
			Token{"2", NumberLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"2", NumberLit, Position{}},

			// vector
			Token{"1", NumberLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"", SlashLit, Position{}},
			Token{"1", NumberLit, Position{}},

			Token{"", Eof, Position{}},
		},
	},
}

type PrintState struct{}

func (ps *PrintState) State(p *Parser) {
	print("STATE: ", p.String(), "\n")
}

func (ps *PrintState) Emit(t *Token) {
	print("EMIT: ", t.String(), "\n")
}
