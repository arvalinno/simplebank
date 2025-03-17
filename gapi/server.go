package gapi

import (
	"fmt"

	"github.com/arvalinno/simplebank/api"
	db "github.com/arvalinno/simplebank/db/sqlc"
	"github.com/arvalinno/simplebank/pb"
	"github.com/arvalinno/simplebank/token"
	"github.com/arvalinno/simplebank/util"
	"github.com/go-playground/validator/v10"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	validate := validator.New()
	validate.RegisterValidation("currency", api.ValidCurrency)

	fmt.Println("Sini")

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
