package tree_sitter_tx

import (
	"testing"

	tree_sitter "github.com/smacker/go-tree-sitter"
)

func TestCanLoadGrammar(t *testing.T) {
	language := tree_sitter.NewLanguage(Language())
	if language == nil {
		t.Errorf("Error loading Tx grammar")
	}
}
