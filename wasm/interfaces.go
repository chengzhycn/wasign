package wasm

import "time"

type Signer interface {
	sign(input SignInput) (SignOutput, error)
}

type SignInput struct {
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Url       string            `json:"url"`
	Body      []byte            `json:"body"`
	AccessKey string            `json:"access_key,omitempty"`
	SecretKey string            `json:"secret_key,omitempty"`
	Token     string            `json:"token,omitempty"`
	IssuedAt  time.Time         `json:"issued_at,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

type SignOutput struct {
	AdditionalHeaders map[string]string `json:"additional_headers,omitempty"`
}
