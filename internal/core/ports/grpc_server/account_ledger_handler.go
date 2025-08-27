package grpc

import (
	"context"

	pb "github.com/ShristiRnr/Finance_mierp/api/pb"
	"github.com/ShristiRnr/Finance_mierp/internal/core/services"
	"github.com/ShristiRnr/Finance_mierp/internal/core/domain"
)

type AccountHandler struct {
	pb.UnimplementedAccountServiceServer
	svc *services.AccountService
}

func NewAccountHandler(svc *services.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
	acc, err := h.svc.CreateAccount(ctx, domain.Account{
		Code:   req.Code,
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AccountResponse{Id: acc.ID.String(), Code: acc.Code, Name: acc.Name}, nil
}
