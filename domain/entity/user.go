package entity

type User struct {
	ID      string `dynamodbav:"id" json:"id"`
	Name    string `dynamodbav:"name" json:"name"`
	Surname string `dynamodbav:"surname" json:"surname"`
}
