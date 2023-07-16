package backends

import (
	"errors"

	"config-wrapper/app/config/confd/backends/content"
	"config-wrapper/app/config/confd/backends/file"
)

// The StoreClient interface is implemented by objects that can retrieve
// key/value pairs from a backend store.
type StoreClient interface {
	GetValues() (map[string]string, error)
}

// New is used to create a storage client based on our configuration.
func New(config Config) (StoreClient, error) {
	switch config.Backend {
	case "file":
		return file.NewFileClient(config.YAMLFile), nil
	case "content":
		return content.NewKubeClient(config.Content), nil
	}

	return nil, errors.New("Invalid backend")
}