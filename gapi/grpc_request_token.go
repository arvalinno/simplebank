package gapi

import (
	"context"

	"github.com/arvalinno/simplebank/pb"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) RequestToken(ctx context.Context, req *pb.RequestTokenRequest) (*pb.RequestTokenResponse, error) {
	_, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderId,
			GrossAmt: req.GetGrossAmt(),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.GetFirstName(),
			Email: req.GetEmail(),
		},
	}

	resp, errMidtrans := snap.CreateTransactionToken(snapReq)
	if errMidtrans != nil {
		return nil, status.Errorf(codes.Internal, "error: %s", errMidtrans.GetMessage())
	}

	grpcResp := &pb.RequestTokenResponse{
		Token: resp,
	}

	return grpcResp, nil
}
