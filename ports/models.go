package ports

type Event struct {
	Source     string
	Detail     string
	DetailType string
	Resources  []string
}
