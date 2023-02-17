package main

import (
	"fmt"
	"os"
	"strings"

	u "utils"

	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
)

// IMPLANT AVAILABLE COMMANDS
func availableCommands() []u.Command {
	exitCmd := u.Command{
		Name:        "exit",
		Description: "Exit SYNTH",
		Help:        "Exit SYNTH",
		SubArgs:     nil,
	}

	helpCmd := u.Command{
		Name:        "help",
		Description: "Show all the available commands",
		Help:        "Show all the available commands",
		SubArgs:     nil,
	}

	setCmd := u.Command{
		Name:        "listener",
		Description: "Start a new listener",
		Help:        "Start a new listener by specifying a host, port, uri, and password",
		SubArgs:     nil,
	}

	showCmd := u.Command{
		Name:        "show",
		Description: "Show listeners|implants",
		Help:        "Show listeners|implants",
		SubArgs: map[string][]string{
			"listeners": nil,
			"implants":  nil,
		},
	}

	// return all command dataclasses as a list
	return []u.Command{exitCmd, helpCmd, setCmd, showCmd}
}

// AUTOCOMPLETE SUGGESTIONS
func completer(d prompt.Document) []prompt.Suggest {
	commands := availableCommands()
	suggestions := make([]prompt.Suggest, 0)

	// Get the current text before the cursor
	inputText := d.TextBeforeCursor()

	// If the input text starts with "show ", suggest only "listeners" and "implants"
	if strings.HasPrefix(inputText, "show ") {
		showSuggestions := []prompt.Suggest{
			{Text: "listeners", Description: "Show all listeners"},
			{Text: "implants", Description: "Show all implants"},
		}
		return prompt.FilterHasPrefix(showSuggestions, d.GetWordAfterCursor(), true)
	}

	// TODO 1: IF INPUT TEXT STARTS WITH listener ALLOW FOR MULTIPLE ARGUMENTS AFTER
	if strings.HasPrefix(inputText, "listener ") {
		showSuggestions := []prompt.Suggest{
			{Text: "<host> <port> <base uri> <password>", Description: "Usage"},
		}
		return prompt.FilterHasPrefix(showSuggestions, d.GetWordAfterCursor(), true)
	}

	// If the input text starts with anything else, suggest available commands
	for _, cmd := range commands {
		if strings.HasPrefix(cmd.Name, d.GetWordBeforeCursor()) {
			suggestions = append(suggestions, prompt.Suggest{Text: cmd.Name, Description: cmd.Description})
		}
	}
	return suggestions
}

func runSYNTH() {
	for {
		input := prompt.Input("[SYNTH] > ", completer)
		if input == "exit" {
			fmt.Println("[-] Exiting SYNTH C2")
			return
		}
		if input == "help" {
			commands := availableCommands()
			// print out all available commands
			for _, cmd := range commands {
				fmt.Printf("%s - %s\n", cmd.Name, cmd.Description)
			}
		}
		if strings.HasPrefix(input, "listener") {
			// code for handling listener command goes here
			createListener(input)

			// add to listeners table
			listenersTable := tablewriter.NewWriter(os.Stdout)
			listenersTable.SetHeader([]string{"Name", "Host", "Port", "Uri", "Password"})
			for _, l := range listeners {
				listenersTable.Append([]string{l.Name, l.Address, l.Port, l.Uri, l.Password})

				// TODO: proceed with prompt after creating the goroutine
				startServer(l.Address, l.Port, l.Uri, l.Password)
			}
			// DEBUG
			// listenersTable.Render()

			// TODO - improvements: ADD ALL THESE TO A DB

		}
		if input == "show" {
			// DEBUGfmt.Println("[*] Choose Listeners OR Implants")
		}
		if input == "show listeners" {
			// display listeners
			listenersTable := tablewriter.NewWriter(os.Stdout)
			listenersTable.SetHeader([]string{"Name", "Host", "Port", "Uri", "Password"})
			for _, l := range listeners {

				listenersTable.Append([]string{l.Name, l.Address, l.Port, l.Uri, l.Password})

			}
			listenersTable.Render()
		}
		if input == "show implants" {
			var implants = []string{"implant1", "implant2", "implant3"} // hardcoded implants

			// display implants
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Type", "Name"})
			for _, i := range implants {

				// TODO: CHANGE HARDCODED TYPE
				table.Append([]string{"x64", i})
			}
			table.Render()
		}

	}

}
