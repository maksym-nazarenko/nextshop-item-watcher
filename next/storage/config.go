package storage

type Config struct {
	Driver  string
	Options map[string]interface{}
}

type NoOptions struct{}
