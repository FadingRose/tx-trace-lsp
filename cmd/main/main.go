package main

import (
	"context"
	"fadingrose/tx-trace-lsp/lsp"
	"fadingrose/tx-trace-lsp/parser"
	"os"

	"go.lsp.dev/jsonrpc2"
)

func main() {
	server := lsp.NewServer()
	ctx := context.Background()
	stream := jsonrpc2.NewStream(os.Stdin)
	conn := jsonrpc2.NewConn(stream)
	server.Run(ctx, conn)
	parser.NewParser() // Initialize the parser to load the grammar
}
