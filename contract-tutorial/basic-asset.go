package main

// Owner structure contains the name of the Company and the name of the Lessor of the basic machine
//type Owner struct {
//	Company string `json:"company"`
//	Lessor  string `json:"lessor"`
//}

// BasicMachine structure contains all the features of a basic machine
type BasicMachine struct {
	ID              string `json:"id"`
	Lessor          string `json:"lessor"`
	Status          string `json:"status"`
	ReservePrice    uint64 `json:"reserveprice"`
	WorkedHours     uint64 `json:"workedhours"`
	PricePerHour    uint64 `json:"priceperhour"`
	Lessee          string `json:"lessee"`
	RentalTime      string `json:"rentaltime"`
	PlaceOfDelivery string `json:"placeofdelivery"`
	WorkHours       uint64 `json:"workhours"`
}

// IotaPayload structure contains the features of an iotapayload message for MAM
type IotaPayload struct {
	Seed     string `json:"seed"`
	MamState string `json:"mamState"`
	Root     string `json:"root"`
	Mode     string `json:"mode"`
	SideKey  string `json:"sideKey"`
}

// SetStatusAvailable function set the status of the machine to available
func (ba *BasicMachine) SetStatusAvailable() {
	ba.Status = "AVAILABLE"
}

// SetStatusReserved function set the status of the machine to reserved
func (ba *BasicMachine) SetStatusReserved() {
	ba.Status = "RESERVED"
}

// SetStatusSent function set the status of the machine to sent
func (ba *BasicMachine) SetStatusSent() {
	ba.Status = "SENT"
}

// SetStatusReceived function set the status of the machine to received
func (ba *BasicMachine) SetStatusReceived() {
	ba.Status = "RECEIVED"
}

// SetStatusWorking function set the status of the machine to working
func (ba *BasicMachine) SetStatusWorking() {
	ba.Status = "WORKING"
}

// SetStatusReturned function set the status of the machine to returned
func (ba *BasicMachine) SetStatusReturned() {
	ba.Status = "RETURNED"
}

// SetStatusInConpany function set the status of the machine to in company
func (ba *BasicMachine) SetStatusInConpany() {
	ba.Status = "IN COMPANY"
}

// SetStatusInMaintenance function set the status of the machine to in maintenance
func (ba *BasicMachine) SetStatusInMaintenance() {
	ba.Status = "IN MAINTENANCE"
}

// SetLessee function set the Lessee of the machine to NO LESSEEE
func (ba *BasicMachine) SetLessee() {
	ba.Lessee = "NO LESSEE"
}

// SetRentalTime function set the rental time of the machine to NO RESERVED TIME
func (ba *BasicMachine) SetRentalTime() {
	ba.RentalTime = "NO RESERVE TIME"
}

// SetPlaceOfDelivery function set the place where te machine will be delivered to NO DESTINATION
func (ba *BasicMachine) SetPlaceOfDelivery() {
	ba.PlaceOfDelivery = "NO DESTINATION"
}

// SetWorkHours function set the work hours of the machine to zero
func (ba *BasicMachine) SetWorkHours() {
	ba.WorkHours = 0
}
