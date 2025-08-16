package hostfuncs

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"

	extism "github.com/extism/go-sdk"
)

// func hmac256(key []byte, toSignString string) ([]byte, error)
func Hmac256(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
	key, err := p.ReadBytes(stack[0])
	if err != nil {
		return
	}

	toSignString, err := p.ReadString(stack[1])
	if err != nil {
		return
	}

	// instantiate hmac
	h := hmac.New(sha256.New, key)
	// write toSignString to hmac
	_, err = h.Write([]byte(toSignString))
	if err != nil {
		return
	}

	// calculate signature and push the memory address to the stack
	signature := h.Sum(nil)
	// TODO: how to handle error in host functions?
	stack[0], _ = p.WriteBytes(signature)
}
