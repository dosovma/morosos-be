package types

type Agreement struct {
	Id        string `dynamodbav:"id" json:"id"`
	Name      string `dynamodbav:"name" json:"name"`
	Surname   string `dynamodbav:"surname" json:"surname"`
	ElapsedAt string `dynamodbav:"elapsed_at" json:"elapsed_at"`
}

type AgreementRange struct {
	Agreements []Agreement `json:"agreements"`
	Next       *string     `json:"next,omitempty"`
}
