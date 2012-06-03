package wfobj

import (
	"fmt"
	"strconv"
)

type meshLoader struct {
	mesh     *Mesh
	vertices VertexList
	tokens   <-chan Token
	token    Token
}

type MeshLoadError string

func NewMeshLoadError(val interface{}) MeshLoadError {
	return MeshLoadError(fmt.Sprintf("%v", val))
}
func (m MeshLoadError) Error() string {
	return "MeshLoadError: " + string(m)
}

// Read a new token from the parser
func (m *meshLoader) read() (ok bool) {
	m.token, ok = <-m.tokens
	return
}

// Read a new token from the parser.
// If the token kind is different from k panic
func (m *meshLoader) readKind(k Kind) {
	m.read()
	m.ensureKind(&m.token, k)
}

func (m *meshLoader) ensureKind(t *Token, k Kind) {
	if t.Kind != k {
		panic(fmt.Sprintf("Invalid token %v. Expecting %v", t, k))
	}
}

// Read a number from the token stream and return the number
// panic if it's not a number
//
// If t is not nil, instead of consuming a token from the stream
// just ensure that t is a valid number literal
func (m *meshLoader) readNumberLit(t *Token) (num float64) {
	if t == nil {
		m.read()
	}
	t = &m.token
	m.ensureKind(t, NumberLit)

	num, err := strconv.ParseFloat(t.Val, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func (m *meshLoader) Load() (err error) {

	defer func() {
		if p := recover(); p != nil {
			err = NewMeshLoadError(p)
		}
	}()

	if !m.read() {
		panic("Empty token stream")
	}

	m.vertices = make(VertexList, 0)
	m.mesh = &Mesh{}
	m.mesh.Faces = make([]Face, 0)

	for {
		switch m.token.Kind {
		case VertexDecl:
			v := Vertex{}
			v.X = float32(m.readNumberLit(nil))
			v.Y = float32(m.readNumberLit(nil))
			v.Z = float32(m.readNumberLit(nil))
			m.vertices = append(m.vertices, v)
		case FaceDecl:
			f := Face{}
			f.Vertices = make(VertexList, 0)
			faceDef := true
			for faceDef {
				if !m.read() {
					break
				}
				switch m.token.Kind {
				case NumberLit:
					idx := int32(m.readNumberLit(&m.token))
					// copy the vertices from the vertex list to the face
					f.Vertices = append(f.Vertices, m.vertices[idx-1])
				default:
					// face definition completed
					faceDef = false
				}
			}
			m.mesh.Faces = append(m.mesh.Faces, f)
			// keep the last token and continue from here
			// prevent the final call to m.read
			// since the previous for consumed all tokens from the channel
			continue
		case Eof:
			break
		default:
			panic(fmt.Sprintf("Unexpected token (%v) expecting: %v", &m.token, fmt.Sprintf("[%v]", []Kind{VertexDecl, FaceDecl, Eof})))
		}

		// advance to the next token
		if !m.read() {
			break
		}
	}

	return
}

// Load a new mesh
func LoadMesh(tokens <-chan Token) (m *Mesh, err error) {
	ml := &meshLoader{nil, nil, tokens, Token{}}
	err = ml.Load()
	m = ml.mesh
	return
}
