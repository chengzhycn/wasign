package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
)

func hmac256(key []byte, toSignString string) ([]byte, error) {
	// 实例化HMAC-SHA256哈希
	h := hmac.New(sha256.New, key)
	// 写入待签名的字符串
	_, err := h.Write([]byte(toSignString))
	if err != nil {
		return nil, err
	}
	// 计算签名并返回
	return h.Sum(nil), nil
}

func main() {
	key := []byte("1234567890")
	toSignString := "Hello, World!"
	signature, err := hmac256(key, toSignString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("signature: %x\n", signature)
}
