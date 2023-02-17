package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ARGUMENTS TO START THE SERVER
var rootCmd = &cobra.Command{
	Use:   "SYNTH",
	Short: "C2 by 0xSH4RKS",
	Long:  "Command and Control framework written by 0xSH4RKS",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the SYNTH Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Starting the SYNTH Server...")
		runSYNTH()
	},
}

var payloadCmd = &cobra.Command{
	Use:   "payload",
	Short: "Generate a SYNTH Payload",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[+] Generating a SYNTH Payload...")
	},
}

// STARTING THE SERVER
func init() {
	rootCmd.AddCommand(serverCmd, payloadCmd)
}

// MAIN LOGIC
func main() {

	fmt.Println("████████████████████████████████████████████")
	fmt.Println("█─▄▄▄▄█▄─█─▄█▄─▀█▄─▄█─▄─▄─█─█─███─▄▄▄─█▀▄▄▀█")
	fmt.Println("█▄▄▄▄─██▄─▄███─█▄▀─████─███─▄─███─███▀██▀▄██")
	fmt.Println("▀▄▄▄▄▄▀▀▄▄▄▀▀▄▄▄▀▀▄▄▀▀▄▄▄▀▀▄▀▄▀▀▀▄▄▄▄▄▀▄▄▄▄▀")
	fmt.Println("             Author: 0xSH4RKS")
	fmt.Println("")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

}
