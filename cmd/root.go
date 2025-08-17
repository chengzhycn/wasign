/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	extism "github.com/extism/go-sdk"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"

	hostfuncs "github.com/chengzhycn/wasign/pkg/host_funcs"
	"github.com/chengzhycn/wasign/wasm"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wasign",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		manifest := extism.Manifest{
			Wasm: []extism.Wasm{
				extism.WasmFile{
					Path: "wasm/jwt/jwt.wasm",
				},
			},
		}

		// a new host function example.
		hmac256 := extism.NewHostFunctionWithStack("hmac256",
			hostfuncs.Hmac256,
			[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR},
			[]extism.ValueType{extism.ValueTypePTR})

		ctx := context.Background()
		config := extism.PluginConfig{
			EnableWasi: true,
		}

		// using compiled plugin for efficiency.
		// Step 1: Compile the plugin once
		compiledPlugin, err := extism.NewCompiledPlugin(ctx, manifest, config, []extism.HostFunction{hmac256})
		if err != nil {
			panic(err)
		}

		// Example: Using the compiled plugin in multiple goroutines
		var wg sync.WaitGroup
		for i := 0; i < 1; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				// Step 2: Create an instance from the compiled plugin
				plugin, err := compiledPlugin.Instance(ctx, extism.PluginInstanceConfig{})
				if err != nil {
					fmt.Printf("Failed to initialize plugin: %v\n", err)
					return
				}
				defer plugin.Close(ctx)

				input := wasm.SignInput{
					AccessKey: "1234567890",
					SecretKey: "1234567890",
					Headers:   map[string]string{"X-Account-UUID": "1234567891", "X-Account-Info": "1234567890"},
					IssuedAt:  time.Now(),
				}

				data, err := json.Marshal(input)
				if err != nil {
					fmt.Println(err)
					return
				}

				_, out, err := plugin.Call("sign", data)
				if err != nil {
					fmt.Println(err)
					return
				}

				var output wasm.SignOutput
				if err := json.Unmarshal(out, &output); err != nil {
					fmt.Println(err)
					return
				}

				// Extract the token from "Bearer <token>" format
				authHeader := output.AdditionalHeaders["Authorization"]
				var tokenString string
				if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
					tokenString = authHeader[7:]
				} else {
					tokenString = authHeader
				}

				// Parse with time validation disabled to handle WASM timing issues
				// WASM environments may not have access to correct system time
				token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
					return []byte(input.SecretKey), nil
				}, jwt.WithValidMethods([]string{"HS256"}))
				if err != nil {
					fmt.Printf("Goroutine %d JWT parse error: %v\n", id, err)
					return
				}

				// Validate token signature and extract claims
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					exp := claims["exp"].(float64)
					fmt.Printf("Goroutine %d - Token signature valid. Claims: app_key=%s, account_uuid=%s, exp=%s\n",
						id, claims["app_key"], claims["account_uuid"], time.Unix(int64(exp), 0).Format(time.RFC3339))
				} else {
					fmt.Printf("Goroutine %d - Failed to extract claims\n", id)
					return
				}

				fmt.Printf("Goroutine %d result: %v\n", id, output)
			}(i)
		}
		wg.Wait()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wasign.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
