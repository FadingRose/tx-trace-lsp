package parser

import (
	"context"
	tree_sitter_tx "fadingrose/tx-trace-lsp/bindings/go"

	sitter "github.com/smacker/go-tree-sitter"
)

// Parser 封装了我们的 Tree-sitter 解析器
type Parser struct {
	*sitter.Parser
}

func NewParser() *Parser {
	p := sitter.NewParser()

	language := sitter.NewLanguage(tree_sitter_tx.Language())

	p.SetLanguage(language)

	return &Parser{p}
}

func (p *Parser) Parse(content []byte) *sitter.Tree {
	tree, _ := p.Parser.ParseCtx(context.Background(), nil, content)
	return tree
}
