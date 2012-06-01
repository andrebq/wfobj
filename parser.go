package wfobj

import (
	"os"
	"io/ioutil"
	"fmt"
	"unicode/utf8"
)

type Parser struct {
	Contents []byte
	VList VertexList
	sz int
	C rune
	pos int
}

func NewParser(fileName string) (p *Parser, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	
	p = &Parser{nil,make(VertexList,0), 0, 0, 0}
	p.Contents, err = ioutil.ReadAll(file)
	return
}

func (p *Parser) Parse() (err error) {
	defer func() {
		var ok bool
		err, ok = recover().(error);
		if !ok {
			err = nil
		}
	}()
	
	for p.Next() {
		switch(p.C) {
			case 'v':
				p.ReadVector()
			case utf8.RuneError:
				panic(fmt.Sprintf("Invalid utf-8 code @ %v", p.pos))
		}
	}
	
	return
}

// Read the x y z[ w] information for a vector
func (p *Parser) ReadVector() {
	//TODO implement this
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
	return true
}

// Push the last run back in the reader
func (p *Parser) PushBack() {
	p.pos -= p.sz
	p.C = utf8.RuneError
}
