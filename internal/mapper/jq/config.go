package jq

import (
	"errors"
	"fmt"
	"time"

	"github.com/leotech-dev/leoflow/internal/common"
)

const PREFIX = "JQ"

const DEFAULT_TIMEOUT = time.Second

const (
	ModeMap = "map"
	ModeTag = "tag"
)

type jqConfig struct {
	Debug      bool
	Mode       string
	Expression string
	Timeout    time.Duration
}

func initConfig(c *jqConfig) error {
	c.Mode = ModeMap
	c.Timeout = DEFAULT_TIMEOUT

	err := common.InitEnvconfig[jqConfig](PREFIX, c)
	if err != nil {
		return err
	}

	if !validModes[c.Mode] {
		return fmt.Errorf("unknown mode: %s", c.Mode)
	}

	if c.Expression == "" {
		return errors.New("empty expression")
	}

	return nil
}

var validModes = map[string]bool{
	ModeMap: true,
	ModeTag: true,
}
