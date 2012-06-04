package wfobj

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

type Kind int
const (
	VertexDecl = Kind(iota)
	FaceDecl
	NumberLit
	Eof
)

var kindNames = map[Kind]string{
	VertexDecl: "VECTOR_DECLARATION",
	FaceDecl:   "FACE_DECLARATION",
	NumberLit:  "NUMBER_LITERAL",
	Eof:        "EOF",
}

func (k Kind) String() string {
	return kindNames[k]
}

type Token struct {
	Val  string
	Kind Kind
	Pos  Position
}

func (t *Token) String() string {
	return fmt.Sprintf("[%v @ %v]%v", t.Kind, &t.Pos, t.Val)
}

type Position struct {
	// line in the stream
	Line int
	// column of the current line
	Col int
}

func (p *Position) String() string {
	return fmt.Sprintf("(line: %v, col: %v)", p.Line, p.Col)
}

type Debug interface {
	State(p *Parser)
	Emit(t *Token)
}

type Parser struct {
	Contents string
	VList    VertexList
	Tokens   chan Token
	Debug    Debug
	sz       int
	C        rune
	// position in the stream
	pos int
	// current line position
	cPos Position
	// old line position
	oPos Position
}

// An error that happened during the parse of the file
type ParseError string

// Return the messsage with the current position of the parser
func NewParseError(p *Parser, msg string) ParseError {
	return ParseError(fmt.Sprintf("%v %v", msg, p.cPos))
}

// Error interface
func (p ParseError) Error() string {
	return string(p)
}

func NewParserFromFile(fileName string) (p *Parser, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	buff, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	p = NewLiteralParser(string(buff))
	return
}

// Parse the contents of the string variable
func NewLiteralParser(literal string) (p *Parser) {
	literal = strings.Replace(literal, "\r\n", "\n", -1)
	p = &Parser{literal, make(VertexList, 0), make(chan Token, 0), nil, 0, 0, 0, Position{1, 1}, Position{1, 0}}
	return
}

// Start the parser and emit the tokens in the Tokens channel
func (p *Parser) Parse() (err error) {
	defer func() {
		close(p.Tokens)
		if val := recover(); val != nil {
			err = NewParseError(p, fmt.Sprintf("%v", val))
		}
	}()

	for p.Next() {
		switch p.C {
		case 'v':
			p.Emit("", VertexDecl)
			p.ReadNumberList()
		case 'f':
			p.Emit("", FaceDecl)
			p.ReadNumberList()
		case '#':
			// comment
			p.DiscardUntil("\n")

		case utf8.RuneError:
			panic(fmt.Sprintf("Invalid utf-8 code @ %v", p.pos))
		}
	}
	p.Emit("", Eof)

	return
}

// Emit a token
func (p *Parser) Emit(val string, kind Kind) {
	t := Token{val, kind, p.cPos}
	if p.Debug != nil {
		p.Debug.Emit(&t)
	}
	p.Tokens <- t
}

// Discard all chars from the stream that match at least one of the chars passed
func (p *Parser) Discard(chars string) {
	for p.NextIf(chars) {
	}
}

// Discard all the runes until the one of the chars is found
func (p *Parser) DiscardUntil(chars string) {
	for p.Next() {
		if strings.IndexAny(string(p.C), chars) != -1 {
			p.PushBack()
			return
		}
	}
}

// Accumulate the runes from the stream while it matches the chars
func (p *Parser) Acc(chars string) string {
	acc := ""
	for p.NextIf(chars) {
		acc += string(p.C)
	}
	return acc
}

// Read a variable length list o numbers
func (p *Parser) ReadNumberList() {
	p.Discard(" ")
	for p.NextIf("0123456789-.") {
		// push the last digit/signal back in the stream
		p.PushBack()
		p.ReadNumberLit()
		p.Discard(" ")
	}
}

// Read the x y z[ w] information for a vector
func (p *Parser) ReadNumberLit() {
	val := ""

	if p.NextIf("-") {
		val += "-"
	}

	val += p.ReadInt()
	if p.NextIf(".") {
		val += "."
		val += p.ReadInt()
	}

	p.Emit(val, NumberLit)
}

// Read a integer and panic if none is found
func (p *Parser) ReadInt() string {
	num := p.Acc("0123456789")
	if len(num) == 0 {
		panic("Expecting one of: 0123456789")
	}
	return num
}

// Check if there is more runes in the contents
func (p *Parser) HasNext() bool {
	return p.pos < len(p.Contents)
}

// Read the rune and move to the next
func (p *Parser) Next() bool {

	// EOF
	if !p.HasNext() {
		return false
	}
	p.C, p.sz = utf8.DecodeRuneInString(p.Contents[p.pos:])
	if p.C == utf8.RuneError {
		return false
	}
	p.pos += p.sz
	// if it is a new line
	// increment the line number
	if p.C == '\n' {
		p.oPos = p.cPos
		p.cPos = Position{p.oPos.Line + 1, 1}
	}
	p.cPos.Col += 1
	return true
}

// Read the rune only if it's in the chars
func (p *Parser) NextIf(chars string) bool {
	ok, _ := p.Peek(chars)
	if ok {
		ok = p.Next()
	}
	return ok
}

// Peek the next run without consuming it
func (p *Parser) Peek(chars string) (ok bool, r rune) {
	ok = true
	r = utf8.RuneError
	if !p.HasNext() {
		ok = false
		return
	}

	r, _ = utf8.DecodeRuneInString(p.Contents[p.pos:])
	if r == utf8.RuneError {
		ok = false
		return
	}

	// no need to check nothing more
	if len(chars) == 0 {
		return
	}

	if strings.IndexAny(string(r), chars) == -1 {
		ok = false
	}
	return
}

// Push the last run back in the reader
func (p *Parser) PushBack() {
	if p.sz == 0 {
		panic("Cannot push more than one time")
	}
	p.pos -= p.sz
	p.sz = 0
	// if it was a new line
	// decrement the line count
	if p.C == '\n' {
		p.cPos = p.oPos
		p.oPos = Position{}
	}
	p.C = utf8.RuneError
}

// Return a string representation of the current state of the parser
func (p *Parser) String() string {
	part := p.Contents[p.pos:]
	if len(part) > 10 {
		part = part[:10]
	}
	return fmt.Sprintf("Contents: %q... @ %v", part, p.cPos)
}
