package types

//go:generate mockgen -destination=./mocks/mock_store.go -package=mocks github.com/aws-samples/serverless-go-demo/types Store

import (
	"context"
)

type Store interface {
	All(context.Context, *string) (AgreementRange, error)
	Get(context.Context, string) (*Agreement, error)
	Put(context.Context, Agreement) error
	Delete(context.Context, string) error
}
