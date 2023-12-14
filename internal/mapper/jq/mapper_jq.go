package jq

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/itchyny/gojq"
	"github.com/leotech-dev/leoflow/pkg/jq"
	nfmapper "github.com/numaproj/numaflow-go/pkg/mapper"
)

type JqMapper struct {
	cfg    jqConfig
	jqFn   jq.Jq
	logger slog.Logger
}

func (m *JqMapper) Init() error {
	err := initConfig(&m.cfg)
	if err != nil {
		return err
	}

	err = m.initExpression()
	if err != nil {
		return err
	}

	return nil
}

func (m *JqMapper) initExpression() error {
	jqFn, err := jq.New(m.cfg.Expression)
	if err != nil {
		return err
	}

	m.jqFn = jqFn

	return nil
}

func (m *JqMapper) Map(ctx context.Context, keys []string, datum nfmapper.Datum) nfmapper.Messages {
	input, err := m.parseInput(datum.Value())
	if err != nil {
		return nfmapper.MessagesBuilder().Append(nfmapper.MessageToDrop())
	}

	ctx, cancel := context.WithTimeout(ctx, m.cfg.Timeout) // in some cases gojq can get stuck with infinite results, hence the limit
	defer cancel()

	return m.processResults(
		m.jqFn(ctx, input),
		keys,
		datum,
	)
}

func (m *JqMapper) processResults(r gojq.Iter, keys []string, datum nfmapper.Datum) nfmapper.Messages {
	switch m.cfg.Mode {
	case ModeTag:
		return m.doTag(r, keys, datum)
	}

	return m.doMap(r, keys)
}

func (m *JqMapper) doMap(r gojq.Iter, keys []string) nfmapper.Messages {
	msgs := nfmapper.MessagesBuilder()
	for {
		v, ok := r.Next()
		if !ok {
			break
		}

		if _, ok := v.(error); ok {
			return msgs.Append(nfmapper.MessageToDrop())
		}

		out, err := json.Marshal(v)
		if err != nil {
			return nfmapper.MessagesBuilder().Append(nfmapper.MessageToDrop())
		}

		msgs = msgs.Append(nfmapper.NewMessage(out).WithKeys(keys))
	}

	return msgs
}

func (m *JqMapper) doTag(r gojq.Iter, keys []string, datum nfmapper.Datum) nfmapper.Messages {
	msgs := nfmapper.MessagesBuilder()
	tags := []string{}
	for {
		v, ok := r.Next()
		if !ok {
			break
		}

		if _, ok := v.(error); ok {
			return msgs.Append(nfmapper.MessageToDrop())
		}

		t, ok := v.(string)
		if !ok {
			return msgs.Append(nfmapper.MessageToDrop())
		}

		tags = append(tags, t)
	}

	return msgs.Append(nfmapper.NewMessage(datum.Value()).WithKeys(keys).WithTags(tags))
}

func (m *JqMapper) parseInput(data []byte) (any, error) {
	var input any
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()

	err := d.Decode(&input)

	return input, err
}
