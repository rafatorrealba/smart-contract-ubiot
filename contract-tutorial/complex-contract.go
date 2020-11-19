package main

// Importing modules
import (
	"fmt"
	"time"
	"errors"
	"encoding/json"

	. "github.com/iotaledger/iota.go/api"
    "github.com/iotaledger/iota.go/trinary"
    "github.com/iotaledger/iota.go/bundle"

	"github.com/rafatorrealba/hlf-iota-conector/iota" 			// Module of IOTA Connector
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ComplexContract contract for handling BasicMachines
type ComplexContract struct {
	contractapi.Contract
}

//Constants and variables
var walletSeed1 string = "UVZULUWZMKQLJISSR99GV9HJGEHXNLYRSOCRLNEIWBLYKIDXMYAGFNXUVZTPMELNHMZOUFPK9AFLOJLUV"
var walletSeed2 string
var walletAddress1 string
var walletAddress2 string = "SHMUHXOSTFSEA9QZEIANRGW9UGXOICU9DIVKOCDIIQPJ9JPUTSHHMRZKPGXUOWSEAKLOTDSGUHBGSWREXMIMWJRNZC"
var walletKeyIndex uint64 = 3
var amount uint64 = 0

// NewMachine function adds a new basic machine to the world state using id as key
func (cc *ComplexContract) NewMachine(ctx CustomTransactionContextInterface, id string, lessor string, reserveprice uint64, workedhours uint64, priceperhour uint64) error {
	existing := ctx.GetData()

	if existing != nil {
		return fmt.Errorf("Cannot create new basic machine in world state as key %s already exists", id)
	}

	// Assinging all features of a new Machine
	ba := new(BasicMachine)
	ba.ID = id
	ba.Lessor = lessor
	ba.ReservePrice = reserveprice
	ba.PricePerHour = priceperhour
	ba.SetLessee()
	ba.SetRentalTime()
	ba.SetPlaceOfDelivery()
	ba.WorkedHours = workedhours
	ba.SetWorkHours()
	ba.SetStatusAvailable()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	err := ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "New machine:", "Root:", root)  			                    	// Printing logs

	return nil
}

// ReserveMachine function changes the Machine status to Reserved and assings a Lessee
func (cc *ComplexContract) ReserveMachine(ctx CustomTransactionContextInterface, id string, lesseeAdd string, rentaltimeAdd string, placeofdeliveryAdd string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Can't update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)

	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "AVAILABLE" {
		return fmt.Errorf("Can't reserve machine %s, because it's not available", id)
	}

	//IOTA Transfer
	amount = ba.ReservePrice                            // Assinging the reserved price as amount to transfer
	TransferIOTA(walletSeed1, walletAddress2, amount)	// Init transfer
	
	// Updating the state of the machine
	ba.Lessee = lesseeAdd
	ba.SetStatusReserved()
	ba.RentalTime = rentaltimeAdd
	ba.PlaceOfDelivery = placeofdeliveryAdd

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "New reserve:", "Root:", root)  			                    	// Printing logs

	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// SentMachine function changes the status of a basic machine to sent
func (cc *ComplexContract) SentMachine(ctx CustomTransactionContextInterface, id string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Can't update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "RESERVED" {
		return fmt.Errorf("Can't sent machine %s, because it's not reserved", id)
	}

	ba.SetStatusSent()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine sent:", "Root:", root)  			                    // Printing logs


	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// ReceivedMachine function changes the status of a basic machine to recived
func (cc *ComplexContract) ReceivedMachine(ctx CustomTransactionContextInterface, id string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Can't update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "SENT" {
		return fmt.Errorf("Can't received machine %s, because it has not been sent", id)
	}

	ba.SetStatusReceived()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine received:", "Root:", root)  			                // Printing logs


	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// PayPerUse function changes the status of a basic machine to Working
func (cc *ComplexContract) PayPerUse(ctx CustomTransactionContextInterface, id string, workhoursAdd uint64) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Can't update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "RECEIVED" && ba.Status != "WORKING" {
		return fmt.Errorf("Can't be put to work the machine %s, because it has not been received", id)
	}

	//IOTA transfer
	amount = ba.PricePerHour * workhoursAdd            	//Assinging the reserved price as amount to transfer
	TransferIOTA(walletSeed1, walletAddress2, amount)	// Init transfer

	// Updating the state of the machine
	ba.SetStatusWorking()
	ba.WorkHours = workhoursAdd
	ba.WorkedHours += workhoursAdd

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine working:", "Root:", root)  			                    // Printing logs


	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// ReturnedMachine function changes the status of a basic machine to returned
func (cc *ComplexContract) ReturnedMachine(ctx CustomTransactionContextInterface, id string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Cant update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "RECEIVED" && ba.Status != "WORKING" {
		return fmt.Errorf("Can't be returned the machine %s, because it has not been received", id)
	}

	// Updating the state of the machine
	ba.SetStatusReturned()
	ba.WorkHours = 0
	ba.SetRentalTime()
	ba.SetPlaceOfDelivery()
	ba.SetWorkHours()
	ba.SetLessee()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine returned:", "Root:", root)  			                // Printing logs


	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// MachineInCompany function changes the status of a basic machine to in company
func (cc *ComplexContract) MachineInCompany(ctx CustomTransactionContextInterface, id string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Cannot update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "RETURNED" {
		return fmt.Errorf("Can't be put to work the machine %s, because it has not been received", id)
	}

	// Updating the state of the machine
	ba.SetStatusInConpany()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine in company:", "Root:", root)  			                // Printing logs


	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// MachineInMaintenance function changes the status of a basic machine to in maintenance
func (cc *ComplexContract) MachineInMaintenance(ctx CustomTransactionContextInterface, id string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Cannot update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "IN COMPANY" {
		return fmt.Errorf("Can't be put to maintenance the machine %s, because it is not in company", id)
	}

	// Updating the state of the machine
	ba.SetStatusInMaintenance()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine in maintenance:", "Root:", root)  			            // Printing logs


	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// AvailableMachine function changes the status of a basic machine to available
func (cc *ComplexContract) AvailableMachine(ctx CustomTransactionContextInterface, id string) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Cannot update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Handing status errors
	if ba.Status != "IN COMPANY" && ba.Status != "IN MAINTENANCE" {
		return fmt.Errorf("Can't make available machine %s, because it is not in company", id)
	}

	// Updating the state of the machine
	ba.SetStatusAvailable()

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Machine available:", "Root:", root)  			                // Printing logs

	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// UpdateReservePrice function updates the reserveprice of a basic machine
func (cc *ComplexContract) UpdateReservePrice(ctx CustomTransactionContextInterface, id string, reservepriceAdd uint64) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Cannot update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Updating the state of the machine
	ba.ReservePrice = reservepriceAdd

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Price updated:", "Root:", root)  			                    // Printing logs

	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// UpdatePricePerHour function updates the priceperhours of a basic machine
func (cc *ComplexContract) UpdatePricePerHour(ctx CustomTransactionContextInterface, id string, priceperhourAdd uint64) error {
	existing := ctx.GetData()

	if existing == nil {
		return fmt.Errorf("Cannot update machine in world state as key %s does not exist", id)
	}

	ba := new(BasicMachine)

	err := json.Unmarshal(existing, ba)
	if err != nil {
		return fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	// Updating the state of the machine
	ba.PricePerHour = priceperhourAdd

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                  // Get te current date and time                                
	_, root, _ := iota.PublishAndReturnState(
		string(baBytes), false, "", "", iota.MamMode, iota.PadSideKey(iota.MamSideKey))		// Inint function to send MAM
	fmt.Println(timestamp, "Price updated:", "Root:", root)  			                    // Printing logs

	err = ctx.GetStub().PutState(id, []byte(baBytes))

	if err != nil {
		return errors.New("Unable to interact with world state")
	}

	return nil
}

// GetMachine function returns the basic machine with id given from the world state
func (cc *ComplexContract) GetMachine(ctx CustomTransactionContextInterface, id string) (*BasicMachine, error) {
	existing := ctx.GetData()

	if existing == nil {
		return nil, fmt.Errorf("Cannot read world state pair with key %s. Does not exist", id)
	}

	ba := new(BasicMachine)
	err := json.Unmarshal(existing, ba)

	if err != nil {
		return nil, fmt.Errorf("Data retrieved from world state for key %s was not of type BasicMachine", id)
	}

	return ba, nil
}

func TransferIOTA(seed string, recipientAddress string, amount uint64) {
    // Connect to a node
    api, err := ComposeAPI(HTTPClientSettings{URI: "https://nodes.thetangle.org:443"})
    must(err)

    // Define an address to which to send IOTA tokens 
    address := trinary.Trytes(recipientAddress)

    // Define an input transaction object
    // that sends 1 i to your new address
    transfers := bundle.Transfers{
          {
              Address: address,
              Value: amount,
          },
      }

    fmt.Println("Sending 1 i to " + walletAddress2);

    trytes, err := api.PrepareTransfers(seed, transfers, PrepareTransfersOptions{})
    must(err)
    
    myBundle, err := api.SendTrytes(trytes, 3, 14)
    must(err)

    hash, _ := json.Marshal(myBundle)
    fmt.Println("HASH: " + string(hash[9:92]) + "\n")
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}

// GetEvaluateTransactions returns functions of ComplexContract not to be tagged as submit
func (cc *ComplexContract) GetEvaluateTransactions() []string {
	return []string{"GetMachine"}
}
