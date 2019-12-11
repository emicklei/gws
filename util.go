package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tj/go-spin"
	"github.com/urfave/cli"
)

// https://developers.google.com/admin-sdk/directory/v1/reference/groups/list
// As an account administrator, you can also use the my_customer alias to represent your account's customerId.
const myAccoutsCustomerID = "my_customer"

func ifZero(i, j int) int {
	if i == 0 {
		return j
	}
	return i
}

func showSpinnerWhile(c *cli.Context) func() {
	// no spinner while verbose logging
	if c.GlobalBool("v") {
		return func() {}
	}
	spinner := spin.New()
	spinner.Set(spin.Box1)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
			default:
				// reprint new spinner state
				fmt.Fprintf(os.Stderr, "\r%s", spinner.Next())
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return func() {
		done <- true
		// remove spinner
		fmt.Fprintf(os.Stderr, "\033[%dD", 1)
	}
}

func optionJSON(c *cli.Context, v interface{}) bool {
	wantsJSON := c.Bool("json")
	if wantsJSON {
		data, _ := json.MarshalIndent(v, "", "\t")
		fmt.Println(string(data))
	}
	return wantsJSON
}

func promptForYes(c *cli.Context, message string) bool {

	// Don't prompt for confirmation if the quiet flag is enabled
	if c.GlobalBool("quiet") {
		return true
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	yn, _ := reader.ReadString('\n')
	return strings.HasPrefix(yn, "Y") || strings.HasPrefix(yn, "y")
}
