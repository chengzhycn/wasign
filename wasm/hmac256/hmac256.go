package main

import pdk "github.com/extism/go-pdk"

//go:wasmimport extism:host/user hmac256
func hmac256(uint64, uint64) uint64

//go:wasmexport hmac256_demo
func hmac256_demo() int32 {
	var input struct {
		Key          string `json:"key"`
		ToSignString string `json:"to_sign_string"`
	}

	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return 1
	}

	key := []byte(input.Key)
	mem1 := pdk.AllocateBytes(key)
	defer mem1.Free()

	toSignString := input.ToSignString
	mem2 := pdk.AllocateString(toSignString)
	defer mem2.Free()

	ptr := hmac256(mem1.Offset(), mem2.Offset())
	rmem := pdk.FindMemory(ptr)
	pdk.Output(rmem.ReadBytes())
	return 0
}

func main() {}
