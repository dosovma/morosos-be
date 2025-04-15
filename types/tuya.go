package types

type TuyaClient interface {
	PostDevice(string) error
}
