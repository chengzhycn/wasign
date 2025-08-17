package main

import (
	"time"

	"github.com/extism/go-pdk"
	"github.com/golang-jwt/jwt/v5"

	"github.com/chengzhycn/wasign/wasm"
)

// func Sign(input SignInput) (SignOutput, error)
//
//go:wasmexport sign
func Sign() int32 {
	var input wasm.SignInput
	if err := pdk.InputJSON(&input); err != nil {
		pdk.SetError(err)
		return 1
	}

	res, err := sign(input)
	if err != nil {
		pdk.SetError(err)
		return 1
	}

	var output wasm.SignOutput
	output.AdditionalHeaders = res.AdditionalHeaders

	if err := pdk.OutputJSON(output); err != nil {
		pdk.SetError(err)
		return 1
	}

	return 0
}

func sign(input wasm.SignInput) (wasm.SignOutput, error) {
	accountUUID := input.Headers["X-Account-UUID"]
	accountInfo := input.Headers["X-Account-Info"]

	// using time.Now() in the wasm plugin will always return 2022-01-01T00:00:00Z.
	// don't know why, but it's a fact. So we use input.IssuedAt as a workaround.
	auth := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"app_key":      input.AccessKey,
		"account_uuid": accountUUID,
		"account_info": accountInfo,
		"exp":          jwt.NewNumericDate(input.IssuedAt.Add(time.Hour * 1)),
	})

	token, err := auth.SignedString([]byte(input.SecretKey))
	if err != nil {
		return wasm.SignOutput{}, err
	}

	return wasm.SignOutput{
		AdditionalHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
	}, nil
}

func main() {}
