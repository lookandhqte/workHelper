package app

import (
	"context"
	"log"
	"net"

	accountAPI "git.amocrm.ru/gelzhuravleva/amocrm_golang/api/grpc/account"
	accountCR "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/controller/grpc/account"
	accountUC "git.amocrm.ru/gelzhuravleva/amocrm_golang/internal/usecase/account"
	grpc "google.golang.org/grpc"
)

func (a *App) StartGRPCServer(uc *accountUC.UseCase) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	accountServer := accountCR.NewServer(uc)
	accountAPI.RegisterAccountServiceServer(s, accountServer)

	a.AddTask(func(ctx context.Context) {
		log.Printf("gRPC server starting on %s", lis.Addr().String())
		if err := s.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	})

	a.AddTask(func(ctx context.Context) {
		<-ctx.Done()
		s.GracefulStop()
		log.Println("gRPC server stopped gracefully")
	})
}
