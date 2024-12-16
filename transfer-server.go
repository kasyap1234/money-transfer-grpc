package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/kasyap1234/money-transfer-grpc/proto"

	"log"

	"google.golang.org/grpc"
)


type Account struct {
	ID string 
	Owner string 
	AccountType string 
	Balance float64

}

type MoneyTransferServer struct {
	pb.UnimplementedMoneyTransferServiceServer
	accounts map[string]*Account

}
func NewMoneyTransferServer()* MoneyTransferServer{
	return &MoneyTransferServer{
		accounts: make(map[string]*Account),
	}
}

func (s *MoneyTransferServer) ValidateTransfer(ctx context.Context, req *pb.ValidationRequest) (*pb.TransferResponse, error) {
	

	sourceAccount, sourceExists := s.accounts[req.SourceAccountId]
	_, destExists := s.accounts[req.DestinationAccountId]

	if !sourceExists || !destExists {
		return &pb.TransferResponse{
			Status:  pb.TransferResponse_ACCOUNT_NOT_FOUND,
			Message: "Source or destination account not found",
		}, nil
	}

	if req.Amount <= 0 {
		return &pb.TransferResponse{
			Status:  pb.TransferResponse_INVALID_AMOUNT,
			Message: "Transfer amount must be positive",
		}, nil
	}

	if sourceAccount.Balance < req.Amount {
		return &pb.TransferResponse{
			Status:  pb.TransferResponse_INSUFFICIENT_FUNDS,
			Message: "Insufficient funds for transfer",
		}, nil
	}

	return &pb.TransferResponse{
		Status:  pb.TransferResponse_SUCCESSFUL,
		Message: "Transfer validated successfully",
	}, nil
}

// first check validate transfer is yes then proceed to execute transfer ; 
func (s *MoneyTransferServer) ExecuteTransfer(ctx context.Context, req *pb.TransferRequest)(*pb.TransferResponse,error) {
    // First validate the transfer
    validationResponse, err := s.ValidateTransfer(ctx, &pb.ValidationRequest{
        SourceAccountId: req.SourceAccountId,
        DestinationAccountId: req.DestinationAccountId,
        Amount: req.Amount,
    })
    if err != nil {
        return nil, err
    }
    
    if validationResponse.Status != pb.TransferResponse_SUCCESSFUL {
        return &pb.TransferResponse{
            Status: validationResponse.Status,
            Message: validationResponse.Message,
        }, nil
    }

   

    // Get accounts
    sourceAccount := s.accounts[req.SourceAccountId]
    destAccount := s.accounts[req.DestinationAccountId]

    // Execute transfer
    sourceAccount.Balance -= req.Amount
    destAccount.Balance += req.Amount

    // Prepare successful response
    return &pb.TransferResponse{
        Status: pb.TransferResponse_SUCCESSFUL,
        Message: "Transfer completed successfully",
        NewSourceBalance: sourceAccount.Balance,
        NewDestinationBalance: destAccount.Balance,
        TransactionId: req.TransactionId,
    }, nil
}

func (s *MoneyTransferServer) GetAccountBalance(ctx context.Context,req *pb.Account)(*pb.Account,error){
	account,exists :=s.accounts[req.AccountId]
	if !exists{
		return nil,fmt.Errorf("account not found")
	}
	return &pb.Account{
		AccountId: account.ID,
		OwnerName: account.Owner,
		AccountType: account.AccountType,
		Balance: account.Balance,

	},nil 

}


func main(){
	lis,err := net.Listen("tcp",":50051")
	
	if err != nil{
		log.Fatalf("failed to listen: %v",err)
	}
	grpcServer := grpc.NewServer()
	moneyServer := NewMoneyTransferServer()
	moneyServer.accounts["1"] = &Account{
		ID: "1",
		Owner: "John Doe",
		AccountType: "Savings",
		Balance: 1000.0,
	}
	moneyServer.accounts["2"] = &Account{
		ID: "2",
		Owner: "Jane Smith",
		AccountType: "Checking",
		Balance: 500.0, 
	}
	moneyServer.accounts["3"]=&Account{
		ID: "3",
		Owner: "Alice Johnson",
		AccountType: "Savings",
		Balance: 2000.0,
	}
	moneyServer.accounts["4"]=&Account{
		ID: "4",
		Owner: "Bob Williams",
		AccountType: "Checking",
		Balance: 1500.0,
	}
	pb.RegisterMoneyTransferServiceServer(grpcServer,moneyServer)
	log.Printf("server listening at %v",lis.Addr())
	if err := grpcServer.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v",err)
	}

}