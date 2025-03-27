package gapi

import (
	"context"
	"fmt"

	"github.com/arvalinno/simplebank/pb"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func (server *Server) RequestToken(ctx context.Context, req *pb.RequestTokenRequest) (*pb.RequestTokenResponse, error) {
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

	resp, err := snap.CreateTransactionToken(snapReq)
	if err != nil {
		fmt.Println("Error :", err.GetMessage())
	}
	fmt.Println("Response : ", resp)

	grpcResp := &pb.RequestTokenResponse{
		Token: resp,
	}

	return grpcResp, nil
}
