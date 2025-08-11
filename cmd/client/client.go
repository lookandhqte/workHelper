package main

import (
	"context"
	"log"

	"git.amocrm.ru/gelzhuravleva/amocrm_golang/api/grpc/account"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatal("Usage: client <account_id>")
	// }
	// accountID, err := strconv.Atoi(os.Args[1])
	// if err != nil {
	// 	log.Fatalf("Invalid account_id: %v", err)
	// }
	accountID := 1

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials())) //NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := account.NewAccountServiceClient(conn)
	resp, err := c.UnsubscribeAccount(context.Background(), &account.UnsubscribeRequest{
		AccountId: int32(accountID),
	})
	if err != nil {
		log.Fatalf("could not unsubscribe: %v", err)
	}
	log.Printf("Response: Success: %v, Message: %s", resp.Success, resp.Message)
}
