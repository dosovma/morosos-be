package entity

type Apartment struct {
	ID       string   `dynamodbav:"id" json:"id"`
	Landlord User     `dynamodbav:"landlord" json:"landlord"`
	Address  string   `dynamodbav:"address" json:"address"`
	Devices  []Device `dynamodbav:"devices" json:"devices"`
}

type Device struct {
	ID   string `dynamodbav:"id" json:"id"`
	Name string `dynamodbav:"name" json:"name"`
	IsOn bool   `dynamodbav:"on" json:"on"`
}

type ApartmentAction = string

const (
	ApartmentOff ApartmentAction = "apartment_off"
	ApartmentOn  ApartmentAction = "apartment_on"
)
