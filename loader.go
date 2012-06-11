package wfobj

import (
	"fmt"
	"strconv"
)

type meshLoader struct {
	mesh     *Mesh
	vertices VertexList
	normals  VertexList
	tokens   []*Token
	pos      int
}

type MeshLoadError string

func NewMeshLoadError(val interface{}) MeshLoadError {
	return MeshLoadError(fmt.Sprintf("%v", val))
}
func (m MeshLoadError) Error() string {
	return "MeshLoadError: " + string(m)
}

// Read a new token from the parser
func (m *meshLoader) next() (ok bool) {
	m.pos++
	if m.pos >= len(m.tokens) {
		return
	}
	ok = true
	return
}

func (m *meshLoader) peek(k Kind) (t *Token, ok bool) {
	npos := m.pos + 1
	if npos >= len(m.tokens) {
		return
	}
	t = m.tokens[npos]
	if k != AnyKind {
		ok = t.Kind == k
	}
	return
}

// Read the current token
func (m *meshLoader) token() *Token {
	if m.pos < 0 {
		panic("Invalid position. Must be non zero")
	}
	if m.pos >= len(m.tokens) {
		panic("Invalid position. Must be less then length")
	}
	return m.tokens[m.pos]
}

func (m *meshLoader) ensureKind(k Kind) {
	if m.token().Kind != k {
		panic(fmt.Sprintf("Invalid token %v. Expecting %v", m.token(), k))
	}
}

func (m *meshLoader) pushBack() {
	m.pos--
}

// Read a number from the token stream and return the number
// panic if it's not a number
//
// If t is not nil, instead of consuming a token from the stream
// just ensure that t is a valid number literal
func (m *meshLoader) readNumberLit() (num float64) {
	m.next()
	t := m.token()
	m.ensureKind(NumberLit)

	num, err := strconv.ParseFloat(t.Val, 64)
	if err != nil {
		panic(err)
	}
	return num
}

// Read the face declaration with the number/number/number format
func (m *meshLoader) readFaceDecl(f *Face) {
	for m.next() {
		t := m.token()
		if t.Kind == NumberLit {
			m.pushBack()
			idx := int32(m.readNumberLit())
			f.Vertices = append(f.Vertices, m.vertices[idx-1])

			// texture information
			m.next()
			t = m.token()
			if t.Kind == SlashLit {
				// TODO handle textures
			} else {
				m.pushBack()
				continue
			}

			// normal information
			m.next()
			t = m.token()
			if t.Kind == SlashLit {
				idx := int32(m.readNumberLit())
				f.Normals = append(f.Normals, m.normals[idx-1])
			} else {
				m.pushBack()
				continue
			}
		} else {
			m.pushBack()
			break
		}
	}
}

func (m *meshLoader) Load() (err error) {

	defer func() {
		if p := recover(); p != nil {
			err = NewMeshLoadError(p)
		}
	}()

	m.vertices = make(VertexList, 0)
	m.normals = make(VertexList, 0)
	m.mesh = &Mesh{}
	m.mesh.Faces = make([]Face, 0)

	for m.next() {
		switch m.token().Kind {
		case VertexDecl:
			v := Vertex{}
			v.X = float32(m.readNumberLit())
			v.Y = float32(m.readNumberLit())
			v.Z = float32(m.readNumberLit())
			m.vertices = append(m.vertices, v)
		case NormalDecl:
			n := Vertex{}
			n.X = float32(m.readNumberLit())
			n.Y = float32(m.readNumberLit())
			n.Z = float32(m.readNumberLit())
			m.normals = append(m.normals, n)
		case FaceDecl:
			f := Face{}
			f.Vertices = make(VertexList, 0)
			f.Normals = make(VertexList, 0)
			m.readFaceDecl(&f)
			m.mesh.Faces = append(m.mesh.Faces, f)
		case Eof:
			break
		default:
			panic(fmt.Sprintf("Unexpected token (%v) expecting: %v", m.token(), fmt.Sprintf("[%v]", []Kind{VertexDecl, FaceDecl, Eof})))
		}
	}

	return
}

// Load a new mesh
func LoadMesh(tokens <-chan *Token) (m *Mesh, err error) {
	ml := &meshLoader{nil, nil, nil, make([]*Token, 0), -1}
	for t := range tokens {
		ml.tokens = append(ml.tokens, t)
	}
	err = ml.Load()
	m = ml.mesh
	return
}

// Load a new mesh from the given .obj file
func LoadMeshFromFile(file string) (m *Mesh, err error) {
	p, err := NewParserFromFile(file)
	if err != nil {
		return
	}
	go p.Parse()
	m, err = LoadMesh(p.Tokens)
	return
}
