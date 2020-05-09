// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package httpendpoint

import (
	"github.com/elastic/beats/v7/libbeat/common/transport/tlscommon"
)

// Config contains information about httpjson configuration
type config struct {
	TLS             *tlscommon.ServerConfig `config:"ssl"`
	ClientAuth      bool                    `config:"ssl.client_auth"`
	BasicAuth       bool                    `config:"basic_auth"`
	Username        string                  `config:"username"`
	Password        string                  `config:"password"`
	ResponseCode    int                     `config:"response_code" validate:"positive"`
	ResponseBody    string                  `config:"response_body"`
	ResponseHeaders string                  `config:"response_headers"`
	ListenAddress   string                  `config:"listen_address"`
	ListenPort      string                  `config:"listen_port"`
	URL             string                  `config:"url"`
	Prefix          string                  `config:"prefix"`
}

func defaultConfig() config {
	return config{
		BasicAuth:       false,
		Username:        "",
		Password:        "",
		ResponseCode:    200,
		ResponseBody:    `{"message": "success"}`,
		ResponseHeaders: `{"Content-Type": "application/json"}`,
		ListenAddress:   "127.0.0.1",
		ListenPort:      "8000",
		URL:             "/",
		Prefix:          "json",
	}
}
