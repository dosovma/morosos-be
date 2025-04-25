package entity

type AgreementStatus = string

const (
	Draft  AgreementStatus = "draft"
	Signed AgreementStatus = "signed"
)

type AgreementAction = string

const (
	Sign AgreementAction = "sign"
)

type Agreement struct {
	ID        string          `dynamodbav:"id" json:"id"`
	StartAt   string          `dynamodbav:"start_at" json:"start_at"`
	ElapsedAt string          `dynamodbav:"elapsed_at" json:"elapsed_at"`
	UpdatedAt string          `dynamodbav:"updated_at" json:"updated_at"`
	Tenant    User            `dynamodbav:"tenant" json:"tenant"`
	Apartment ApartmentData   `dynamodbav:"apartment" json:"apartment"`
	Status    AgreementStatus `dynamodbav:"status" json:"status"`
	Text      string          `dynamodbav:"text" json:"text"` // TODO move out from entity to service layer
}

type ApartmentData struct {
	ApartmentID string `dynamodbav:"apartment_id" json:"apartment_id"`
	Landlord    User   `dynamodbav:"landlord" json:"landlord"`
	Address     string `dynamodbav:"address" json:"address"`
}
