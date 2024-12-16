package main

import (
    "context"
    "log"
    "time"
    pb "github.com/kasyap1234/money-transfer-grpc/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
   
conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    // Create client
    client := pb.NewMoneyTransferServiceClient(conn)

    // Set timeout for API call
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Test GetAccountBalance
    account, err := client.GetAccountBalance(ctx, &pb.Account{AccountId: "1"})
    if err != nil {
        log.Fatalf("could not get balance: %v", err)
    }
    log.Printf("Account Balance: %v", account.Balance)

    // Test ExecuteTransfer
    transferReq := &pb.TransferRequest{
        SourceAccountId: "1",
        DestinationAccountId: "2",
        Amount: 100.0,
        TransactionId: "tx1",
        Description: "Test transfer",
    }

    response, err := client.ExecuteTransfer(ctx, transferReq)
    if err != nil {
        log.Fatalf("Transfer failed: %v", err)
    }

    log.Printf("Transfer Status: %v", response.Status)
    log.Printf("Message: %v", response.Message)
    log.Printf("New Source Balance: %v", response.NewSourceBalance)
    log.Printf("New Destination Balance: %v", response.NewDestinationBalance)
}
