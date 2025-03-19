package gapi

import (
	"context"

	db "github.com/arvalinno/simplebank/db/sqlc"
	"github.com/arvalinno/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ListAccount(ctx context.Context, req *pb.ListAccountRequest) (*pb.ListAccountResponse, error) {

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
	}

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.GetPageSize(),
		Offset: (req.GetPageId() - 1) * req.GetPageSize(),
	}

	dbAccounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error get list accounts: %s", err)
	}

	pbAccounts := []*pb.Account{}
	for i := range dbAccounts {
		tempAccount := &pb.Account{
			Id:      dbAccounts[i].ID,
			Owner:   dbAccounts[i].Owner,
			Balance: dbAccounts[i].Balance,
		}
		pbAccounts = append(pbAccounts, tempAccount)
	}

	rsp := &pb.ListAccountResponse{
		Accounts: pbAccounts,
	}

	return rsp, nil
}
