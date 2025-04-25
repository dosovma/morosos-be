package entity

type TemplateName = string

const (
	AgreementTemplateName TemplateName = "agreement"
)

type Template struct {
	Name string `dynamodbav:"name" json:"name"`
	Text string `dynamodbav:"text" json:"text"`
}
