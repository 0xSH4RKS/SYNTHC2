package main

import (
	"net/http"
	u "utils"

	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
)

var listeners []u.Listener // declare an empty list of Listeners

func createListener(userInput string) {
	usage := "[-] Usage: listener <host> <port> <base uri> <password>"
	args := strings.Fields(userInput)

	if len(args) < 5 {
		fmt.Println("[-] Something in your arguments is wrong, fix it!")
		fmt.Println(usage)
		return
	}

	host := args[1]
	port := args[2]
	baseURI := args[3]
	password := args[4]
	fmt.Printf("[+] Starting listener on %s:%s%s with password %s\n", host, port, baseURI, password)

	// create a new listener
	listener := u.Listener{
		Name:     "listener1",
		Address:  host,
		Port:     port,
		Uri:      baseURI,
		Password: password,
	}

	listeners = append(listeners, listener)

	// TODO: using martini:
	// 1. start thread -> goroutines
	// 2. start_server()
	// 3. opens handler on arguments from Listener Data Class
}

func startServer(host string, port string, baseURI string, password string) {
	// Set the application environment to "production"
	envy.Set("GO_ENV", "production")

	// Create a new Buffalo application
	app := buffalo.New(buffalo.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Prefix:       baseURI,
		SessionName:  "buffalo-session",
		SessionStore: nil,
	})

	// Define a handler function for the server
	app.GET("/", func(c buffalo.Context) error {
		c.Set("name", "Thomas")

		return c.Render(http.StatusOK, render.String("Hi <%= name %>"))
	})

	// Start the server in a goroutine
	go func() {
		// Start the server
		fmt.Printf("[+] Starting server at http://%s%s\n", app.Addr, baseURI)
		if err := app.Serve(); err != nil {
			fmt.Printf("[-] Error starting server: %v\n", err)
		}
	}()

	// Check if a password was provided
	if password != "" {
		fmt.Printf("[+] Password for server: %s\n", password)
	}
}
