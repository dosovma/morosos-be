package types

type Agreement struct {
	ID          string `dynamodbav:"id" json:"id"`
	ElapsedAt   string `dynamodbav:"elapsed_at" json:"elapsed_at"`
	Tenant      User   `dynamodbav:"tenant" json:"tenant"`
	ApartmentID string `dynamodbav:"apartment" json:"apartment"`
}

type Apartment struct {
	ID       string   `dynamodbav:"id" json:"id"`
	Landlord User     `dynamodbav:"landlord" json:"landlord"`
	Devices  []Device `dynamodbav:"devices" json:"devices"`
}

type User struct {
	Name    string `dynamodbav:"name" json:"name"`
	Surname string `dynamodbav:"surname" json:"surname"`
}

type Device struct {
	ID   string `dynamodbav:"id" json:"id"`
	Name string `dynamodbav:"name" json:"name"`
	IsOn bool   `dynamodbav:"on" json:"on"`
}
