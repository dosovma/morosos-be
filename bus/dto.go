package bus

type EventType = string

const (
	Agreement EventType = "agreement"
)

type Event struct {
	Source     string
	Detail     string
	DetailType string
	Resources  []string
}
