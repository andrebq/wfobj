package wfobj

import (
	"testing"
)

func discard(ch <-chan Token, done chan bool, t *testing.T) {
	for tok := range ch {
		t.Logf("Token: %v", tok)
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
