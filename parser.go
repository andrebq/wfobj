package wfobj

import (
	"os"
	"io/ioutil"
	"fmt"
	"unicode/utf8"
	"strings"
)

type Kind int

const (
	VectorDecl = Kind(iota)
	FaceDecl
	NumberLit
)

type Token struct {
	Val	string
	Kind Kind
}

type Position struct {
	// line in the stream
	Line int
	// column of the current line
	Col int
}

type Parser struct {
	Contents []byte
	VList VertexList
	Tokens	chan Token
	sz int
	C rune
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
func NewParseError(p *Parser, msg string) ParseError{
	return ParseError(fmt.Sprintf("%v (line: %v, col: %v)", msg, p.cPos.Line, p.cPos.Col))
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
	p = newLiteralParser(string(buff))
	return
}

// Parse the contents of the string variable
func newLiteralParser(literal string) (p *Parser) {
	literal = strings.Replace(literal, "\r\n", "\n", -1)
	p = &Parser{[]byte(literal),make(VertexList,0), make(chan Token,1), 0, 0, 0, Position{1,0}, Position{1,0}}
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
		switch(p.C) {
			case 'v':
				p.Tokens <- Token{"vector", VectorDecl}
				p.ReadNumberLit()
			case utf8.RuneError:
				panic(fmt.Sprintf("Invalid utf-8 code @ %v", p.pos))
		}
	}
	
	return
}

// Discard all chars from the stream that match at least one of the chars passed
func (p *Parser) Discard(chars string) {
	for p.NextIf(chars) {	}
}

// Accumulate the runes from the stream while it matches the chars
func (p *Parser) Acc(chars string) string {
	acc := ""
	for p.NextIf(chars) { acc += string(p.C) }
	return acc
}

// Read the x y z[ w] information for a vector
func (p *Parser) ReadNumberLit() {
	p.Discard(" ")
	tok := Token{"", NumberLit}
	
	tok.Val = p.ReadInt()
	if p.NextIf(".") {
		tok.Val += "."
	}
	tok.Val += p.ReadInt()
	
	p.Tokens <- tok
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
	if !p.HasNext() { return false }
	p.C, p.sz = utf8.DecodeRune(p.Contents[p.pos:])
	if p.C == utf8.RuneError {
		return false
	}
	p.pos += p.sz
	// if it is a new line
	// increment the line number
	if p.C == '\n' {
		p.oPos = p.cPos
		p.cPos = Position{p.oPos.Line + 1, 0}
	}
	p.cPos.Col += 1
	return true
}

// Read the rune only if it's in the chars
func (p *Parser) NextIf(chars string) bool {
	ret := p.Next()
	if ret {
		if strings.IndexAny(string(p.C), chars) == -1 {
			p.PushBack()
			ret = false
		}
	}
	return ret
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
