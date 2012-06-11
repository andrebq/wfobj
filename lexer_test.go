package wfobj

import (
	"testing"
)

func discard(ch <-chan *Token, done chan bool, t *testing.T) {
	idx := 0
	for tok := range ch {
		if tok.Kind != tokens[idx].Kind {
			t.Errorf("Expecting %v at position %v, but got %v", tokens[idx], idx, tok)
		}
		idx++
	}
	done <- true
}

func TestParser(t *testing.T) {
	p := NewLiteralParser(objlit)
	//	p.Debug = &PrintState{}
	done := make(chan bool)
	go discard(p.Tokens, done, t)
	err := p.Parse()
	if err != nil {
		t.Fatalf("Unable to parse the cube.obj file. %v", err)
	}
	<-done
}
