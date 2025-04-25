package entity

type AgreementStatus = string

const (
	Draft  AgreementStatus = "draft"
	Signed AgreementStatus = "signed"
)

type Agreement struct {
	ID        string          `dynamodbav:"id" json:"id"`
	StartAt   string          `dynamodbav:"start_at" json:"start_at"`
	ElapsedAt string          `dynamodbav:"elapsed_at" json:"elapsed_at"`
	UpdatedAt string          `dynamodbav:"updated_at" json:"updated_at"`
	Tenant    User            `dynamodbav:"tenant" json:"tenant"`
	Apartment Apartment       `dynamodbav:"apartment" json:"apartment"`
	Status    AgreementStatus `dynamodbav:"status" json:"status"`
	Text      string          `dynamodbav:"text" json:"text"` // TODO move out from entity to service layer
}

type Apartment struct {
	ID       string   `dynamodbav:"id" json:"id"`
	Landlord User     `dynamodbav:"landlord" json:"landlord"`
	Address  string   `dynamodbav:"address" json:"address"`
	Devices  []Device `dynamodbav:"devices" json:"devices"`
}

type Template struct {
	Name string `dynamodbav:"name" json:"name"`
	Text string `dynamodbav:"text" json:"text"`
}

type User struct {
	ID      string `dynamodbav:"id" json:"id"`
	Name    string `dynamodbav:"name" json:"name"`
	Surname string `dynamodbav:"surname" json:"surname"`
}

type Device struct {
	ID   string `dynamodbav:"id" json:"id"`
	Name string `dynamodbav:"name" json:"name"`
	IsOn bool   `dynamodbav:"on" json:"on"`
}

type Status struct {
	Action Action
}

type Action = string

const (
	Sign         Action = "sign"
	ApartmentOff Action = "apartment_off"
	ApartmentOn  Action = "apartment_on"
)
