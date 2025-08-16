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

	extism "github.com/extism/go-sdk"
	"github.com/spf13/cobra"

	hostfuncs "github.com/chengzhycn/wasign/pkg/host_funcs"
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
					Path: "wasm/hmac256/hmac256_demo.wasm",
				},
			},
		}

		hmac256 := extism.NewHostFunctionWithStack("hmac256",
			hostfuncs.Hmac256,
			[]extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR},
			[]extism.ValueType{extism.ValueTypePTR})

		ctx := context.Background()
		config := extism.PluginConfig{
			EnableWasi: true,
		}

		// Step 1: Compile the plugin once
		compiledPlugin, err := extism.NewCompiledPlugin(ctx, manifest, config, []extism.HostFunction{hmac256})
		if err != nil {
			panic(err)
		}

		// Example: Using the compiled plugin in multiple goroutines
		var wg sync.WaitGroup
		for i := 0; i < 3; i++ {
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

				type Hmac256Input struct {
					Key          string `json:"key"`
					ToSignString string `json:"to_sign_string"`
				}

				data, err := json.Marshal(Hmac256Input{
					Key:          "1234567890",
					ToSignString: "Hello, World!",
				})
				if err != nil {
					fmt.Println(err)
					return
				}

				_, out, err := plugin.Call("hmac256_demo", data)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Printf("Goroutine %d result: %x\n", id, out)
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
