package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"
	admin "google.golang.org/api/admin/directory/v1"
)

func cmdUserList(c *cli.Context) error {

	client := sharedAuthClient()

	srv, err := admin.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve directory Client %v", err)
	}

	r, err := srv.Users.List().
		Customer("my_customer"). // ??
		MaxResults(int64(IfZero(c.Int("limit"), 100))).
		OrderBy("email").Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve users in domain: %v", err)
	}

	wantsJSON := c.String("format") == "JSON"
	if wantsJSON {
		data, _ := json.MarshalIndent(r.Users, "", "\t")
		fmt.Println(string(data))
		return nil
	}
	for _, u := range r.Users {
		// email is default
		fmt.Println(u.PrimaryEmail)
	}
	return nil
}
