package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/arvalinno/simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	fields := strings.Fields(values[0])
	authType := fields[0]
	if strings.ToLower(authType) != authorizationBearer {
		return nil, fmt.Errorf("invalid authorization type: %s", authType)
	}

	authToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(authToken)
	fmt.Println("payload: ", payload)
	if err != nil {
		return nil, fmt.Errorf("token is invalid: %s", err)
	}

	return payload, nil
}
