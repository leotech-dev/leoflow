package jq

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type fetchOpts struct {
	Method  string
	Headers map[string][]string
	Body    any
}

func fetch(pipedInput any, args []any) any {
	url, ok := args[0].(string)
	if !ok || url == "" {
		return errors.New("URL must be specified")
	}

	opts, err := parseFetchOpts(args)
	if err != nil {
		return err
	}

	r := resty.New().R()
	r.SetBody(opts.Body)
	r.SetHeaderMultiValues(opts.Headers)

	resp, err := r.Execute(opts.Method, url)
	if err != nil {
		return err
	}

	return map[string]any{
		"status_code": resp.StatusCode(),
		"headers":     resp.Header(),
		"body":        string(resp.Body()),
	}
}

func parseFetchOpts(args []any) (fetchOpts, error) {
	opts := fetchOpts{
		Method:  http.MethodGet,
		Headers: map[string][]string{},
		Body:    nil,
	}

	if len(args) <= 1 {
		return opts, nil
	}

	optsMap, ok := args[1].(map[string]any)
	if !ok {
		return opts, errors.New("fetch opts parameter must be a JSON object")
	}

	if _, ok := optsMap["headers"]; ok {
		hmap, ok := optsMap["headers"].(map[string]any)
		if !ok {
			return opts, errors.New("headers option must be a JSON object")
		}

		for k, v := range hmap {
			hSingleVal, ok := v.(string)
			if ok {
				opts.Headers[k] = []string{hSingleVal}
				continue
			}

			hMultiVal, ok := v.([]string)
			if ok {
				opts.Headers[k] = hMultiVal
				continue
			}

			return opts, fmt.Errorf("header value must be a string or an array of strings: %s", k)
		}
	}

	if _, ok := optsMap["method"]; ok {
		method, ok := optsMap["method"].(string)
		if !ok {
			return opts, errors.New("method option must be a string")
		}

		opts.Method = method
	}

	if _, ok := optsMap["body"]; ok {
		opts.Body = optsMap["body"]
	}

	return opts, nil
}
