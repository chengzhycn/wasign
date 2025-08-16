package wasm

type Signer interface {
	sign() int32
}

type SignInput struct {
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	Url       string            `json:"url"`
	Body      []byte            `json:"body"`
	AccessKey string            `json:"access_key,omitempty"`
	SecretKey string            `json:"secret_key,omitempty"`
	Token     string            `json:"token,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}
