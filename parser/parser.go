package parser

import (
	"context"
	"fmt"

	sitter_tx "github.com/fadingrose/tree-sitter-tx/bindings/go"
	sitter "github.com/smacker/go-tree-sitter"
)

type Parser struct {
	*sitter.Parser
}

func NewParser() *Parser {
	p := sitter.NewParser()

	p.SetLanguage(sitter.NewLanguage(sitter_tx.Language()))

	_, err := p.ParseCtx(context.Background(), nil, []byte{})
	if err != nil {
		panic(fmt.Sprintf("error initializing parser: %v", err))
	}

	return &Parser{p}
}

func (p *Parser) MustParse(content []byte) *sitter.Tree {
	tree, err := p.ParseCtx(context.Background(), nil, content)
	if err != nil {
		panic(err)
	}
	return tree
}

func (p *Parser) Parse(content []byte) *sitter.Tree {
	tree, _ := p.ParseCtx(context.Background(), nil, content)
	return tree
}
