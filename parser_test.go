package wfobj

import (
	"testing"
)

func discard(ch <-chan Token, t *testing.T) {
	for tok := range ch {
		t.Logf("Token: %v", tok)
	}
}

func TestParser(t *testing.T) {
	p, err := NewParserFromFile("cube.obj")
	if err != nil {
		t.Fatalf("Unable to create parser. %v", err)
	}
	go discard(p.Tokens, t)
	err = p.Parse()
	if err != nil {
		t.Fatalf("Unable to parse the cube.obj file. %v", err)
	}
}
