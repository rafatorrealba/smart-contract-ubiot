package main

// Importing modules
import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/iotaledger/iota.go/address"
	"github.com/iotaledger/iota.go/consts"
	"github.com/rafatorrealba/hlf-iota-conector/iota" // Module of IOTA Connector
)

// ComplexContract contract for handling BasicMachines
type ComplexContract struct {
	contractapi.Contract
}

//Constants and variables
const walletSeed1 = "XP9ATINAUBHOURCEBVXVORWWRMNIVYIMZRSEV9HYWNMADLIERCI9LYRODBDIXBTEEMMVCELHJUATNJJD9"
const walletSeed2 = "FBBNBNBQPGVFLGNPDNJZVYN9CJO9BGUAALXCJMWTDHIDOR9HOOBRHEYANQDVFKWPXWRPNLMTYJ9LCJFLD"

var walletKeyIndex uint64 = 4
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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "New machine:", ba.ID, "Root:", root, mamState, seed)                      // Printing logs

	return nil //Original return.
	//return shim.Success(nil) //New return of IOTA. Actually unused.
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
	walletAddress, err := address.GenerateAddress(walletSeed2, 0, consts.SecurityLevelMedium, true)
	amount = ba.ReservePrice                                                               // Assinging the reserved price as amount to transfer
	newKeyIndex := iota.TransferTokens(walletSeed1, walletKeyIndex, walletAddress, amount) // Init transfer and get new key index
	walletKeyIndex = newKeyIndex                                                           // Assinging the new index to current index

	// Updating the state of the machine
	ba.Lessee = lesseeAdd
	ba.SetStatusReserved()
	ba.RentalTime = rentaltimeAdd
	ba.PlaceOfDelivery = placeofdeliveryAdd

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "New reserve:", ba.ID, "Root:", root, mamState, seed)                      // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine sent:", ba.ID, "Root:", root, mamState, seed)                     // Printing logs

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

	// Inint function to send MAM
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine recived:", ba.ID, "Root:", root, mamState, seed)                  // Printing logs

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
	walletAddress, err := address.GenerateAddress(walletSeed2, 0, consts.SecurityLevelMedium, true)
	amount = ba.PricePerHour * workhoursAdd                                                // Assinging the reserved price as amount to transfer
	newKeyIndex := iota.TransferTokens(walletSeed1, walletKeyIndex, walletAddress, amount) // Init transfer and get new key index
	walletKeyIndex = newKeyIndex                                                           // Assinging the new index to current index

	// Updating the state of the machine
	ba.SetStatusWorking()
	ba.WorkHours = workhoursAdd
	ba.WorkedHours += workhoursAdd

	// Changing GO structure to JSON format
	baBytes, _ := json.MarshalIndent(ba, "", "  ")

	// Communication with IOTA through MAM
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "New work order:", ba.ID, "Root:", root, mamState, seed)                   // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine returned:", ba.ID, "Root:", root, mamState, seed)                 // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine in Company:", ba.ID, "Root:", root, mamState, seed)               // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine in Maintenance:", ba.ID, "Root:", root, mamState, seed)           // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine available:", ba.ID, "Root:", root, mamState, seed)                // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine reserve price updated:", ba.ID, "Root:", root, mamState, seed)    // Printing logs

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
	timestamp := time.Now().String()[0:19]                                                            // Get te current date and time
	mode := iota.MamMode                                                                              //
	sideKey := iota.PadSideKey(iota.MamSideKey)                                                       //
	mamState, root, seed := iota.PublishAndReturnState(string(baBytes), false, "", "", mode, sideKey) // Inint function to send MAM
	mamState = ""                                                                                     // Setted to ""; unused
	seed = ""                                                                                         // Setted to ""; unused
	fmt.Println(timestamp, "Machine price per hour updated:", ba.ID, "Root:", root, mamState, seed)   // Printing logs

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

// GetEvaluateTransactions returns functions of ComplexContract not to be tagged as submit
func (cc *ComplexContract) GetEvaluateTransactions() []string {
	return []string{"GetMachine"}
}
