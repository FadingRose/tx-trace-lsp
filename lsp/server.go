package lsp

import (
	"context"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

type Server struct {
	conn jsonrpc2.Conn
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(ctx context.Context, conn jsonrpc2.Conn) {
	s.conn = conn
}

// Initialize is called when the client initializes the server.
func (s *Server) Initialize(ctx context.Context, params *protocol.InitializeParams) (*protocol.InitializeResult, error) {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			FoldingRangeProvider: true,
		},
	}, nil
}

func (s *Server) FoldingRangeProvider(ctx context.Context, params *protocol.FoldingRangeParams) ([]*protocol.FoldingRange, error) {
	return nil, nil
}
