package gapi

import (
	"fmt"

	db "github.com/giadat1599/small_bank/db/sqlc"
	"github.com/giadat1599/small_bank/pb"
	"github.com/giadat1599/small_bank/token"
	"github.com/giadat1599/small_bank/utils"
)

// Server serves gRPC requests for banking service
type Server struct {
	pb.UnimplementedSmallBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker	
}



// NewServer creates a new gRPC server 
func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	
	return server, nil
}


