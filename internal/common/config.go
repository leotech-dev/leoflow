package common

import (
	"github.com/kelseyhightower/envconfig"
)

func InitEnvconfig[T interface{}](prefix string, cfg *T) error {
	return envconfig.Process(prefix, cfg)
}
