package gapi

import (
	"context"

	"github.com/arvalinno/simplebank/api"
	db "github.com/arvalinno/simplebank/db/sqlc"
	"github.com/arvalinno/simplebank/pb"
	"github.com/arvalinno/simplebank/val"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type contextAuthPayloadKey string

const AuthorizationPayloadKeyGrpc contextAuthPayloadKey = api.AuthorizationPayloadKey

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
	}

	violations := validateCurrencyRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.GetCurrency(),
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pgxErr, ok := err.(*pgconn.PgError); ok {
			switch pgxErr.Code {
			case "23503", "23505":
				return nil, status.Errorf(codes.PermissionDenied, "error db: %s", err)
			}
		}

		return nil, status.Errorf(codes.Internal, "error creating account: %s", err)
	}

	rsp := &pb.CreateAccountResponse{
		Account: &pb.Account{
			Id:        account.ID,
			Owner:     account.Owner,
			Balance:   account.Balance,
			Currency:  account.Currency,
			CreatedAt: timestamppb.New(account.CreatedAt.Time),
		},
	}

	return rsp, nil
}

func validateCurrencyRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateCurrency(req.GetCurrency())
	if err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}

	return violations
}
