package wfobj

import (
	"testing"
)

func discard(ch <-chan *Token, done chan []*Token, t *testing.T) {
	buff := make([]*Token, 0)
	for tok := range ch {
		buff = append(buff, tok)
	}
	done <- buff
}

func TestParser(t *testing.T) {
	for _, test := range testdata {
		if test.ignore { continue }
		t.Logf("Title: %v", test.title)
		p := NewLiteralParser(test.objlit)
		//p.Debug = &PrintState{}
		done := make(chan []*Token)
		go discard(p.Tokens, done, t)
		err := p.Parse()
		if err != nil {
			t.Errorf("Unable to parse the cube.obj file. %v", err)
		}
		tmp := <-done
		if len(tmp) != len(test.tokens) {
			t.Errorf("Expecting %v tokens but got %v", len(test.tokens), len(tmp))
		}
		for i, tok := range tmp {
			if tok.Kind != test.tokens[i].Kind {
				t.Errorf("Expecting %v but got %v", test.tokens[i], tok)
			}
		}
	}
}
