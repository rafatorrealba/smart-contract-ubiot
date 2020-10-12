package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	simpleContract := new(SimpleContract)
	simpleContract.TransactionContextHandler = new(CustomTransactionContext)
	simpleContract.BeforeTransaction = GetWorldState
	simpleContract.UnknownTransaction = UnknownTransactionHandler

	complexContract := new(ComplexContract)
	complexContract.TransactionContextHandler = new(CustomTransactionContext)
	complexContract.BeforeTransaction = GetWorldState

	simpleContract.Name = "org.example.com.SimpleContract"
	complexContract.Name = "org.example.com.ComplexContract"

	cc, err := contractapi.NewChaincode(simpleContract, complexContract)

	if err != nil {
		panic(err.Error())
	}

	if err := cc.Start(); err != nil {
		panic(err.Error())
	}

	cc.DefaultContract = complexContract.GetName()

}
