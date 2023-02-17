package main

import (
	"io"
	"math/rand"
	"os"
	"time"
	u "utils"

	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

var listeners []u.Listener // declare an empty list of Listeners

func createListener(userInput string) {

	// nice colors :)
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)
	green := color.New(color.FgGreen)
	boldGreen := green.Add(color.Bold)
	magenta := color.New(color.FgHiMagenta)

	usage := "Usage: listener <host> <port> <base uri> <password>"
	args := strings.Fields(userInput)

	now := time.Now()

	if len(args) < 5 {
		boldRed.Print("[-] ")
		fmt.Println("Something in your arguments is wrong, fix it!")

		boldRed.Print("[-] ")
		fmt.Println(usage)
		return
	}

	name := generateRandomString(5)
	host := args[1]
	port := args[2]
	baseURI := args[3]
	password := args[4]

	boldGreen.Print("[+]")
	magenta.Printf(" %d/%d/%d, %d:%d:%d", now.Day(), now.Month(), now.Year(), now.Hour(), now.Minute(), now.Second())
	fmt.Printf(" - Listener started: [%s] http://%s:%s%s (%s)\n", name, host, port, baseURI, password)

	// create a new listener
	listener := u.Listener{
		Name:     name,
		Address:  host,
		Port:     port,
		Uri:      baseURI,
		Password: password,
	}

	listeners = append(listeners, listener)
}

func startServer(host string, port string, baseURI string, password string) {
	// colors for prompt
	red := color.New(color.FgRed)
	boldRed := red.Add(color.Bold)

	gin.DisableConsoleColor() // Disable console color

	// Logging to a file.
	f, _ := os.Create("logs/debug.log")
	gin.DefaultWriter = io.MultiWriter(f)

	// Starting a webserver instance using GIN
	router := gin.Default()
	router.GET(baseURI+"/stage/:id", func(c *gin.Context) {
		stageId := c.Param("id")
		c.JSON(200, gin.H{
			"message": "Malware to be downloaded: " + stageId + "!",
		})
	})

	// Start the server in a goroutine
	go func() {
		// Start the server
		// boldGreen.Print("[+]")
		// fmt.Printf(" Starting server at http://%s"+":%s"+"%s\n", host, port, baseURI)
		if err := router.Run(":" + port); err != nil {
			boldRed.Print("[-]")
			fmt.Printf(" Error starting server: %v\n", err)
		}

	}()

	// // Check if a password was provided
	// if password != "" {
	// 	boldGreen.Print("[+]")
	// 	fmt.Printf(" Password for server: %s\n", password)
	// }
}

func generateRandomString(length int) string {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Define the set of characters that can be used in the random string
	letters := "abcdefghijklmnopqrstuvwxyz"

	// Generate a random string of the specified length
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	// Convert the byte slice to a string and return it
	return string(b)
}
