package account

import (
	"context"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/api/grpc/account"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	account.UnimplementedAccountServiceServer
	uc *accountUC.UseCase
}

func NewServer(uc *accountUC.UseCase) *Server {
	return &Server{uc: uc}
}

func (s *Server) UnsubscribeAccount(ctx context.Context, req *account.UnsubscribeRequest) (*account.UnsubscribeResponse, error) {
	accountID := int(req.GetAccountId())

	if err := s.uc.Delete(accountID); err != nil {
		return &account.UnsubscribeResponse{
			Success: false,
			Message: "Failed to delete account: " + err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &account.UnsubscribeResponse{
		Success: true,
		Message: "Account and all related data deleted successfully",
	}, nil
}
