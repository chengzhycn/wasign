package main

import (
	"fmt"

	"github.com/extism/go-pdk"
)

//go:wasmexport add
func add() int32 {
	var input struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return 1
	}

	res := input.A + input.B
	pdk.OutputString(fmt.Sprintf("%d", res))

	return 0
}

func main() {}
