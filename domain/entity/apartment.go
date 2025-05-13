package entity

type ApartmentAction = string

const (
	ApartmentOff ApartmentAction = "apartment_off"
	ApartmentOn  ApartmentAction = "apartment_on"
)

type Apartment struct {
	ID       string   `dynamodbav:"id" json:"id"`
	Landlord User     `dynamodbav:"landlord" json:"landlord"`
	Address  string   `dynamodbav:"address" json:"address"`
	Devices  []Device `dynamodbav:"devices" json:"devices"`
}

type Device struct {
	ID          string  `dynamodbav:"id" json:"id"`
	Name        string  `dynamodbav:"name" json:"name"`
	IsOn        bool    `dynamodbav:"on" json:"on"`
	PhoneNumber *string `dynamodbav:"phone_number" json:"phone_number"`
}

type ApartmentRange struct {
	Apartments []Apartment `json:"apartments"`
	Next       *string     `json:"next,omitempty"`
}
