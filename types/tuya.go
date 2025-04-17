package types

type TuyaClient interface {
	PostDevice(string, bool) error
}
